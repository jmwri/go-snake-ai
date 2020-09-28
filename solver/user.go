package solver

import (
	"github.com/hajimehoshi/ebiten"
	"go-snake-ai/direction"
	"go-snake-ai/state"
)

func NewUserSolver() *UserSolver {
	return &UserSolver{
		lastPressed: direction.None,
		listening:   false,
	}
}

type UserSolver struct {
	lastPressed direction.Direction
	listening   bool
}

func (s *UserSolver) Name() string {
	return "user"
}

func (s *UserSolver) Init() {
	s.lastPressed = direction.None
	if !s.listening {
		s.listening = true
		go s.listen()
	}
}

func (s *UserSolver) listen() {
	for {
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			s.lastPressed = direction.Up
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			s.lastPressed = direction.Right
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			s.lastPressed = direction.Down
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			s.lastPressed = direction.Left
		}
	}
}

func (s *UserSolver) NextMove(st *state.State) direction.Direction {
	last := s.lastPressed
	s.lastPressed = direction.None
	return last
}
