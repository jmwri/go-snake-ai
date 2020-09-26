package input

import (
	"github.com/hajimehoshi/ebiten"
	"go-snake-ai/direction"
	"go-snake-ai/state"
	"time"
)

func NewUserInput() *UserInput {
	return &UserInput{
		lastPressed: direction.None,
		ticks:       20,
		listening:   false,
	}
}

type UserInput struct {
	lastPressed direction.Direction
	ticks       int
	listening   bool
}

func (i *UserInput) Name() string {
	return "user"
}

func (i *UserInput) Init() {
	i.lastPressed = direction.None
	i.ticks = 20
	if !i.listening {
		i.listening = true
		go i.listen()
		go i.listenSpeed()
	}
}

func (i *UserInput) listen() {
	for {
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			i.lastPressed = direction.Up
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			i.lastPressed = direction.Right
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			i.lastPressed = direction.Down
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			i.lastPressed = direction.Left
		}
	}
}

func (i *UserInput) listenSpeed() {
	for {
		if ebiten.IsKeyPressed(ebiten.KeyComma) {
			i.ticks++
		} else if ebiten.IsKeyPressed(ebiten.KeyPeriod) && i.ticks > 0 {
			i.ticks--
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (i *UserInput) Ticks() int {
	return i.ticks
}

func (i *UserInput) NextMove(s *state.State) direction.Direction {
	last := i.lastPressed
	i.lastPressed = direction.None
	return last
}
