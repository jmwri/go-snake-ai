package solver

import (
	"github.com/hajimehoshi/ebiten"
	"go-snake-ai/direction"
	"go-snake-ai/state"
	"time"
)

func NewUserSolver() *UserSolver {
	return &UserSolver{
		lastPressed: direction.None,
		ticks:       20,
		listening:   false,
	}
}

type UserSolver struct {
	lastPressed direction.Direction
	ticks       int
	listening   bool
}

func (s *UserSolver) Name() string {
	return "user"
}

func (s *UserSolver) Init() {
	s.lastPressed = direction.None
	s.ticks = 20
	if !s.listening {
		s.listening = true
		go s.listen()
		go s.listenSpeed()
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

func (s *UserSolver) listenSpeed() {
	for {
		if ebiten.IsKeyPressed(ebiten.KeyComma) {
			s.ticks++
		} else if ebiten.IsKeyPressed(ebiten.KeyPeriod) && s.ticks > 0 {
			s.ticks--
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (s *UserSolver) Ticks() int {
	return s.ticks
}

func (s *UserSolver) NextMove(st *state.State) direction.Direction {
	last := s.lastPressed
	s.lastPressed = direction.None
	return last
}
