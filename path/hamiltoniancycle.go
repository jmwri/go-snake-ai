package path

import (
	"go-snake-ai/state"
	"go-snake-ai/tile"
)

func NewHamiltonianCycle(longest *BreadthFirstSearchLongest) *HamiltonianCycle {
	return &HamiltonianCycle{
		longest: longest,
	}
}

type HamiltonianCycle struct {
	longest *BreadthFirstSearchLongest
}

func (g *HamiltonianCycle) Generate(state *state.State, from *tile.Vector, to *tile.Vector) (Path, bool) {
	// We don't care about creating a path that ends with `to`.
	// We create a path around the whole map that should include `to`.
	potentialEndVectors := g.getPotentialTailVectors(state)
	if len(potentialEndVectors) == 0 {
		return Path{}, false
	}

	for _, endV := range potentialEndVectors {
		path, found := g.longest.Generate(state, from, endV)
		if !found {
			continue
		}
		// The cycle needs to cover ALL tiles
		if len(path) != state.TotalTiles() {
			continue
		}
		// The cycle needs to end on the starting point
		path = append(path, path[0])
		return path, true
	}

	return Path{}, false
}

func (g *HamiltonianCycle) getPotentialTailVectors(state *state.State) []*tile.Vector {
	potentialVectors := make([]*tile.Vector, 0)
	head := state.SnakeHead()
	tail := state.SnakeTail()
	// If snake has some length, the return only the tail vector
	if *head != *tail {
		potentialVectors = append(potentialVectors, tail)
		return potentialVectors
	}

	// We need give a fake tail that can move into the snakes head vector
	// Return all adjacent vectors that are valid, and are not in front of the snake
	snakeDir := state.SnakeDir()
	adjacentVectors := tile.AdjacentVectors(head)

	for dir, adj := range adjacentVectors {
		if dir == snakeDir {
			continue
		}
		if !state.ValidPosition(adj.X, adj.Y) {
			continue
		}
		potentialVectors = append(potentialVectors, adj)
	}

	return potentialVectors
}
