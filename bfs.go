package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

type BreadthFirstSearch struct {
	Frontier []*Node
	Game     *Maze
}

func (bfs *BreadthFirstSearch) GetFrontier() []*Node {
	return bfs.Frontier
}

func (bfs *BreadthFirstSearch) Add(i *Node) {
	bfs.Frontier = append(bfs.Frontier, i)
}

func (bfs *BreadthFirstSearch) ContainsState(i *Node) bool {
	for _, x := range bfs.Frontier {
		if x.State == i.State {
			return true
		}
	}
	return false
}

func (bfs *BreadthFirstSearch) Empty() bool {
	return len(bfs.Frontier) == 0
}

func (bfs *BreadthFirstSearch) Remove() (*Node, error) {
	if len(bfs.Frontier) > 0 {
		if bfs.Game.Debug {
			fmt.Println("Frontier before remove:")
			for _, x := range bfs.Frontier {
				fmt.Println("Node:", x.State)
			}
		}

		node := bfs.Frontier[0]

		bfs.Frontier = bfs.Frontier[1:]
		return node, nil
	}
	return nil, errors.New("frontier is empty")
}

func (bfs *BreadthFirstSearch) Solve() {
	fmt.Println("Starting to solve maze using breadth first search...")
	bfs.Game.NumExplored = 0

	start := Node{
		State:  bfs.Game.Start,
		Parent: nil,
		Action: "",
	}

	bfs.Add(&start)

	bfs.Game.CurrentNode = &start

	for {
		if bfs.Empty() {
			return
		}

		currentNode, err := bfs.Remove()
		if err != nil {
			log.Println(err)
			return
		}

		if bfs.Game.Debug {
			fmt.Println("Removed", currentNode.State)
			fmt.Println("-------")
			fmt.Println()
		}

		bfs.Game.CurrentNode = currentNode
		bfs.Game.NumExplored++

		//Have we found the solution?

		if bfs.Game.Goal == currentNode.State {
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

			bfs.Game.Solution = Solution{
				Actions: actions,
				Cells:   cells,
			}

			bfs.Game.Expolored = append(bfs.Game.Expolored, currentNode.State)
			break
		}

		bfs.Game.Expolored = append(bfs.Game.Expolored, currentNode.State)

		if bfs.Game.Animate {
			bfs.Game.OutPutImage(fmt.Sprintf("tmp/%06d.png", bfs.Game.NumExplored))
		}

		for _, x := range bfs.Neighbors(currentNode) {
			if !bfs.ContainsState(x) {
				if !inExplored(x.State, bfs.Game.Expolored) {
					bfs.Add(&Node{
						State:  x.State,
						Parent: currentNode,
						Action: x.Action,
					})
				}
			}
		}
	}
}

func (bfs *BreadthFirstSearch) Neighbors(node *Node) []*Node {
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
		if 0 <= x.State.Row && x.State.Row < bfs.Game.Height {
			if 0 <= x.State.Col && x.State.Col < bfs.Game.Width {
				if !bfs.Game.Walls[x.State.Row][x.State.Col].wall {
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
