package main

func inExplored(needle Point, haystack []Point) bool {
	for _, x := range haystack {
		if x.Row == needle.Row && x.Col == needle.Col {
			return true
		}
	}

	return false
}
