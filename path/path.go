package path

import (
	"go-snake-ai/direction"
	"go-snake-ai/tile"
)

type Path []*tile.Vector
type Moves []direction.Direction

func PathToMoves(p Path) Moves {
	moves := make(Moves, 0)
	var prev *tile.Vector
	for i, v := range p {
		if i == 0 {
			prev = v
			continue
		}
		dirToV := tile.DirToVector(prev, v)
		prev = v
		if dirToV == direction.None {
			panic("failed to get direction")
		}
		moves = append(moves, dirToV)
	}
	return moves
}
