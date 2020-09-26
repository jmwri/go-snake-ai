package path

import (
	"go-snake-ai/state"
	"go-snake-ai/tile"
)

type Generator interface {
	Generate(state *state.State, from *tile.Vector, to *tile.Vector) (Path, bool)
}
