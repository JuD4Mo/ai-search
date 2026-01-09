package main

import "os"

func inExplored(needle Point, haystack []Point) bool {
	for _, x := range haystack {
		if x.Row == needle.Row && x.Col == needle.Col {
			return true
		}
	}

	return false
}

func emptyTmp() {
	directory := "./tmp/"
	dir, _ := os.Open(directory)
	filesToDelete, _ := dir.Readdir(0)

	for i := range filesToDelete {
		f := filesToDelete[i]
		fullPath := directory + f.Name()
		_ = os.Remove(fullPath)
	}
}
