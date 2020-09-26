package input

import (
	"go-snake-ai/direction"
	"go-snake-ai/state"
)

type Input interface {
	Name() string
	NextMove(s *state.State) direction.Direction
	Init()
}
