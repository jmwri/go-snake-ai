package path

import (
	"fmt"
	"go-snake-ai/state"
	"go-snake-ai/tile"
	"math/rand"
)

func NewBreadthFirstSearchLongest(bfs *BreadthFirstSearch) *BreadthFirstSearchLongest {
	return &BreadthFirstSearchLongest{
		bfs: bfs,
	}
}

type PrefParallelDirection int

const (
	PrefRandom PrefParallelDirection = iota
	PrefNegative
	PrefPositive
)

type BreadthFirstSearchLongest struct {
	bfs          *BreadthFirstSearch
	prefParallel PrefParallelDirection
}

func (g *BreadthFirstSearchLongest) SetPrefParallel(pp PrefParallelDirection) {
	g.prefParallel = pp
}

func (g *BreadthFirstSearchLongest) Generate(state *state.State, from *tile.Vector, to *tile.Vector) (Path, bool) {
	path, found := g.bfs.Generate(state, from, to)
	if !found {
		return path, found
	}
	return g.expandPath(state, path, to), true
}

func (g *BreadthFirstSearchLongest) expandPath(state *state.State, p Path, to *tile.Vector) Path {
	occupiedVectors := make(map[tile.Vector]bool)
	longestPath := make(Path, 0)
	for _, v := range p {
		occupiedVectors[*v] = true
		longestPath = append(longestPath, v)
	}

	i := len(longestPath) - 1
	for i >= 1 {
		a := longestPath[i]
		b := longestPath[i-1]

		parallelVectors, err := g.parallelVectors(a, b)
		if err != nil {
			i--
			continue
		}

		freeParallelVectors := make([][2]*tile.Vector, 0)
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
			freeParallelVectors = append(freeParallelVectors, side)
		}

		extendedPath := false
		if len(freeParallelVectors) > 0 {
			var parallelVector [2]*tile.Vector
			if len(freeParallelVectors) > 1 {
				// If we have a choice of parallel vectors, then check if we prefer one expansion direction over another
				if g.prefParallel == PrefNegative {
					// Negative favours moving left or up
					parallelVector = freeParallelVectors[0]
				} else if g.prefParallel == PrefPositive {
					// Positive favours moving right or down
					parallelVector = freeParallelVectors[1]
				} else if g.prefParallel == PrefRandom {
					parallelVector = freeParallelVectors[rand.Intn(len(freeParallelVectors))]
				}
			} else {
				parallelVector = freeParallelVectors[0]
			}
			ap := parallelVector[0]
			bp := parallelVector[1]

			// Extend the length by 2
			longestPath = append(longestPath, longestPath[len(longestPath)-2:]...)
			// Shift elements after i right by 2 places
			copy(longestPath[i+2:], longestPath[i:len(longestPath)-2])
			// Add in the 2 expanded vectors
			longestPath[i+1] = ap
			longestPath[i] = bp
			occupiedVectors[*ap] = true
			occupiedVectors[*bp] = true

			extendedPath = true
		}

		if extendedPath {
			i += 2
		} else {
			i--
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
