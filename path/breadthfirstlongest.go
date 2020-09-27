package path

import (
	"fmt"
	"go-snake-ai/state"
	"go-snake-ai/tile"
)

func NewBreadthFirstSearchLongest(bfs *BreadthFirstSearch) *BreadthFirstSearchLongest {
	return &BreadthFirstSearchLongest{
		bfs: bfs,
	}
}

type BreadthFirstSearchLongest struct {
	bfs *BreadthFirstSearch
}

func (g *BreadthFirstSearchLongest) Generate(state *state.State, from *tile.Vector, to *tile.Vector) (Path, bool) {
	path, found := g.bfs.Generate(state, from, to)
	if !found {
		return path, found
	}
	longestPath := g.expandPath(state, path, to)
	return longestPath, true
}

func (g *BreadthFirstSearchLongest) expandPath(state *state.State, p Path, to *tile.Vector) Path {
	occupiedVectors := make(map[tile.Vector]bool)
	longestPath := make(Path, 0)
	for _, v := range p {
		occupiedVectors[*v] = true
		longestPath = append(longestPath, v)
	}

	i := 0
	longestPathLen := len(longestPath) - 1
	for i < longestPathLen {
		a := longestPath[i]
		b := longestPath[i+1]

		parallelVectors, err := g.parallelVectors(a, b)
		if err != nil {
			continue
		}

		extendedPath := false
		for _, side := range parallelVectors {
			// a/b parallel
			ap := side[0]
			bp := side[1]

			if !g.canOccupyVector(ap, state, to, occupiedVectors) {
				continue
			}
			if !g.canOccupyVector(bp, state, to, occupiedVectors) {
				continue
			}

			// Extend the length by 2
			longestPath = append(longestPath, longestPath[len(longestPath)-2:]...)
			// Shift elements after i right by 2 places
			copy(longestPath[i+2:], longestPath[i:len(longestPath)-2])
			// Add in the 2 expanded vectors
			longestPath[i+1] = ap
			longestPath[i+2] = bp
			occupiedVectors[*ap] = true
			occupiedVectors[*bp] = true

			extendedPath = true
			longestPathLen = len(longestPath) - 1
			break
		}

		if extendedPath {
			i = 0
		} else {
			i++
		}
	}
	return longestPath
}

func (g *BreadthFirstSearchLongest) parallelVectors(a *tile.Vector, b *tile.Vector) ([2][2]*tile.Vector, error) {
	parallel := [2][2]*tile.Vector{}

	//   pxp
	//   pxp
	if a.X == b.X {
		parallel = [2][2]*tile.Vector{
			{
				tile.NewVector(a.X-1, a.Y),
				tile.NewVector(b.X-1, b.Y),
			},
			{
				tile.NewVector(a.X+1, a.Y),
				tile.NewVector(b.X+1, b.Y),
			},
		}
		return parallel, nil
	}

	//   pp
	//   xx
	//   pp
	if a.Y == b.Y {
		parallel = [2][2]*tile.Vector{
			{
				tile.NewVector(a.X, a.Y-1),
				tile.NewVector(b.X, b.Y-1),
			},
			{
				tile.NewVector(a.X, a.Y+1),
				tile.NewVector(b.X, b.Y+1),
			},
		}
		return parallel, nil
	}

	return parallel, fmt.Errorf("unable to calculate parallel vectors")
}

func (g *BreadthFirstSearchLongest) canOccupyVector(v *tile.Vector, state *state.State, to *tile.Vector, occupiedVectors map[tile.Vector]bool) bool {
	if !state.ValidPosition(v.X, v.Y) {
		return false
	}
	t := state.Tile(v.X, v.Y)
	if t != tile.TypeNone {
		return false
	}
	if _, ok := occupiedVectors[*v]; ok {
		return false
	}
	if *v == *to {
		return false
	}
	return true
}
