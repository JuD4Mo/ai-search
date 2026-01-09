package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

type DepthFirstSearch struct {
	Frontier []*Node
	Game     *Maze
}

func (dfs *DepthFirstSearch) GetFrontier() []*Node {
	return dfs.Frontier
}

func (dfs *DepthFirstSearch) Add(i *Node) {
	dfs.Frontier = append(dfs.Frontier, i)
}

func (dfs *DepthFirstSearch) ContainsState(i *Node) bool {
	for _, x := range dfs.Frontier {
		if x.State == i.State {
			return true
		}
	}
	return false
}

func (dfs *DepthFirstSearch) Empty() bool {
	return len(dfs.Frontier) == 0
}

func (dfs *DepthFirstSearch) Remove() (*Node, error) {
	if len(dfs.Frontier) > 0 {
		if dfs.Game.Debug {
			fmt.Println("Frontier before remove:")
			for _, x := range dfs.Frontier {
				fmt.Println("Node:", x.State)
			}
		}

		node := dfs.Frontier[len(dfs.Frontier)-1]
		dfs.Frontier = dfs.Frontier[:len(dfs.Frontier)-1]

		return node, nil
	}

	return nil, errors.New("frontier is empty")
}

func (dfs *DepthFirstSearch) Solve() {
	fmt.Println("Starting to solve maze using Depth First Search...")

	dfs.Game.NumExplored = 0

	start := Node{
		State:  dfs.Game.Start,
		Parent: nil,
		Action: "",
	}

	dfs.Add(&start)
	dfs.Game.CurrentNode = &start

	for {
		if dfs.Empty() {
			return
		}

		currentNode, err := dfs.Remove()
		if err != nil {
			log.Println(err)
			return
		}

		if dfs.Game.Debug {
			fmt.Println("Removed", currentNode.State)
			fmt.Println("------")
			fmt.Println()
		}

		dfs.Game.CurrentNode = currentNode
		dfs.Game.NumExplored++

		//Have we found the solution
		if dfs.Game.Goal == currentNode.State {
			var actions []string
			var cells []Point

			for {
				if currentNode.Parent != nil {
					actions = append(actions, currentNode.Action)
					cells = append(cells, currentNode.State)
					currentNode = currentNode.Parent
				} else {
					break
				}
			}

			slices.Reverse(actions)
			slices.Reverse(cells)

			dfs.Game.Solution = Solution{
				Actions: actions,
				Cells:   cells,
			}
			dfs.Game.Expolored = append(dfs.Game.Expolored, currentNode.State)
			break
		}

		dfs.Game.Expolored = append(dfs.Game.Expolored, currentNode.State)

		//Build animation frame if appropiate.

		if dfs.Game.Animate {
			dfs.Game.OutPutImage(fmt.Sprintf("tmp/%06d.png", dfs.Game.NumExplored))
		}

		for _, x := range dfs.Neighbors(currentNode) {
			if !dfs.ContainsState(x) {
				if !inExplored(x.State, dfs.Game.Expolored) {
					dfs.Add(&Node{
						State:  x.State,
						Parent: currentNode,
						Action: x.Action,
					})
				}
			}
		}
	}
}

func (dfs *DepthFirstSearch) Neighbors(node *Node) []*Node {
	row := node.State.Row
	col := node.State.Col

	candidates := []*Node{
		{
			State: Point{
				Row: row - 1,
				Col: col,
			},
			Parent: node,
			Action: "up",
		}, {
			State: Point{
				Row: row + 1,
				Col: col,
			},
			Parent: node,
			Action: "down",
		}, {
			State: Point{
				Row: row,
				Col: col - 1,
			},
			Parent: node,
			Action: "left",
		}, {
			State: Point{
				Row: row,
				Col: col + 1,
			},
			Parent: node,
			Action: "right",
		},
	}

	var neighbors []*Node

	for _, x := range candidates {
		if 0 <= x.State.Row && x.State.Row < dfs.Game.Height {
			if 0 <= x.State.Col && x.State.Col < dfs.Game.Width {
				if !dfs.Game.Walls[x.State.Row][x.State.Col].wall {
					neighbors = append(neighbors, x)
				}
			}
		}
	}

	for i := range neighbors {
		j := rand.Intn(i + 1)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}

	return neighbors
}
