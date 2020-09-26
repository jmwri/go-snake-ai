package path

import (
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

		adjVectors := []*tile.Vector{
			// Up
			{
				X: v.X,
				Y: v.Y - 1,
			},
			// Right
			{
				X: v.X + 1,
				Y: v.Y,
			},
			// Down
			{
				X: v.X,
				Y: v.Y + 1,
			},
			// Left
			{
				X: v.X - 1,
				Y: v.Y,
			},
		}

		for _, adj := range adjVectors {
			if state.ValidPosition(adj.X, adj.Y) {
				t := state.Tile(adj.X, adj.Y)
				if t == tile.TypeNone || t == tile.TypeFruit {
					if _, ok := seenVectors[*adj]; !ok {
						seenVectors[*adj] = true
						vectorOwners[*adj] = v
						q.Add(adj)
					}
				}
			}
		}
	}
}
