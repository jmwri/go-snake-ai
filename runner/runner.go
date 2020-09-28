package runner

import (
	"fmt"
	"go-snake-ai/score"
	"go-snake-ai/solver"
	"go-snake-ai/state"
	"math/rand"
	"time"
)

func NewGameRunner(tileNumX int, tileNumY int, slvr solver.Solver, writer score.Writer) *GameRunner {
	return &GameRunner{
		tileNumX: tileNumX,
		tileNumY: tileNumY,
		slvr:     slvr,
		writer:   writer,
	}
}

type GameRunner struct {
	s          *state.State
	ended      bool
	tileNumX   int
	tileNumY   int
	tileWidth  float64
	tileHeight float64
	slvr       solver.Solver
	writer     score.Writer
}

func (r *GameRunner) Init() {
	rand.Seed(time.Now().UnixNano())
	r.s = state.NewState(r.tileNumX, r.tileNumY)
	r.slvr.Init()
	r.ended = false
}

func (r *GameRunner) TileNumX() int {
	return r.tileNumX
}

func (r *GameRunner) TileNumY() int {
	return r.tileNumY
}

func (r *GameRunner) Ended() bool {
	return r.ended
}

func (r *GameRunner) State() *state.State {
	return r.s
}

func (r *GameRunner) Update() error {
	if r.ended {
		return fmt.Errorf("cant update ended game")
	}

	nextDirection := r.slvr.NextMove(r.s)
	alive, err := r.s.Move(nextDirection)
	if err != nil {
		return err
	}
	if !alive {
		r.ended = true
		r.writer.Write(r.s.Score(), r.s.MaxScore(), r.slvr)
	}

	return nil
}
