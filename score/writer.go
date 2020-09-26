package score

import "go-snake-ai/solver"

type Writer interface {
	Write(score int, maxScore int, slvr solver.Solver)
}
