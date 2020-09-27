package path

import (
	"go-snake-ai/direction"
	"go-snake-ai/queue"
	"go-snake-ai/state"
	"go-snake-ai/tile"
)

func NewBreadthFirstSearch() *BreadthFirstSearch {
	return &BreadthFirstSearch{}
}

type BreadthFirstSearch struct {
}

func (g *BreadthFirstSearch) backtrace(v *tile.Vector, owners map[tile.Vector]*tile.Vector) Path {
	q := queue.NewFILO()
	q.Add(v)
	for {
		owner, ok := owners[*v]
		if ok {
			q.Add(owner)
			v = owner
		} else {
			break
		}
	}

	path := Path{}
	for {
		v := q.Pop()
		if v == nil {
			return path
		}
		path = append(path, v)
	}
}

func (g *BreadthFirstSearch) Generate(state *state.State, from *tile.Vector, to *tile.Vector) (Path, bool) {
	seenVectors := make(map[tile.Vector]bool)
	vectorOwners := make(map[tile.Vector]*tile.Vector)
	q := queue.NewFIFO()
	q.Add(from)
	seenVectors[*from] = true

	for {
		v := q.Pop()
		if v == nil {
			return make(Path, 0), false
		}

		if *v == *to {
			return g.backtrace(v, vectorOwners), true
		}

		adjVectors := tile.AdjacentVectors(v)

		for dir, adj := range adjVectors {
			if _, ok := seenVectors[*adj]; ok {
				continue
			}
			if state.ValidPosition(adj.X, adj.Y) {
				if *v == *from {
					// If we're on the first vector, we need to make sure we don't try to move backwards!
					if direction.IsOpposite(state.SnakeDir(), dir) {
						continue
					}
				}
				t := state.Tile(adj.X, adj.Y)
				if t == tile.TypeNone || t == tile.TypeFruit {
					seenVectors[*adj] = true
					vectorOwners[*adj] = v
					q.Add(adj)
				}
			}
		}
	}
}
