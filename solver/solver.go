package solver

import (
	"go-snake-ai/direction"
	"go-snake-ai/state"
)

type Solver interface {
	Name() string
	NextMove(s *state.State) direction.Direction
	Ticks() int
	Init()
}
