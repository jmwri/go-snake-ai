package score

import "go-snake-ai/input"

type Writer interface {
	Write(score int, maxScore int, input input.Input)
}
