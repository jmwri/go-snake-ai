package input

import (
	"github.com/hajimehoshi/ebiten"
	"go-snake-ai/direction"
	"go-snake-ai/state"
)

func NewUserInput() *UserInput {
	return &UserInput{
		lastPressed: direction.None,
		listening:   false,
	}
}

type UserInput struct {
	lastPressed direction.Direction
	listening   bool
}

func (i *UserInput) Init() {
	i.lastPressed = direction.None
	if !i.listening {
		i.listening = true
		go i.listen()
	}
}

func (i *UserInput) listen() {
	for {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			i.lastPressed = direction.Up
		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			i.lastPressed = direction.Right
		} else if ebiten.IsKeyPressed(ebiten.KeyS) {
			i.lastPressed = direction.Down
		} else if ebiten.IsKeyPressed(ebiten.KeyA) {
			i.lastPressed = direction.Left
		}
	}
}

func (i *UserInput) NextMove(s *state.State) direction.Direction {
	last := i.lastPressed
	i.lastPressed = direction.None
	return last
}
