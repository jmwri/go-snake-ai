package input

import (
	"go-snake-ai/direction"
	"go-snake-ai/state"
)

type Input interface {
	NextMove(s *state.State) direction.Direction
	Init()
}
