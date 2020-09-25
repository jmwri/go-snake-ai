package input

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"go-snake-ai/direction"
	"go-snake-ai/state"
)

func NewUserInput() *UserInput {
	input := &UserInput{
		lastPressed: direction.DirectionNone,
	}
	go input.listen()
	return input
}

type UserInput struct {
	lastPressed direction.Direction
}

func (i *UserInput) listen() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		i.lastPressed = direction.DirectionUp
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		i.lastPressed = direction.DirectionRight
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		i.lastPressed = direction.DirectionDown
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		i.lastPressed = direction.DirectionLeft
	}
}

func (i *UserInput) NextMove(s *state.State) direction.Direction {
	last := i.lastPressed
	i.lastPressed = direction.DirectionNone
	return last
}
