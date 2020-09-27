package solver

import (
	"go-snake-ai/direction"
	"go-snake-ai/path"
	"go-snake-ai/state"
)

type RegenType int

const (
	RegenEveryTick RegenType = iota
	RegenEveryFruit
	RegenNever
)

func NewPathFollowingSolver(name string, pathGen path.Generator, regenPath RegenType) *PathFollowingSolver {
	return &PathFollowingSolver{
		name:      name,
		pathGen:   pathGen,
		regenPath: regenPath,
	}
}

type PathFollowingSolver struct {
	name        string
	pathGen     path.Generator
	regenPath   RegenType
	moves       path.Moves
	currentMove int
	prevScore   int
}

func (s *PathFollowingSolver) Name() string {
	return s.name
}

func (s *PathFollowingSolver) Init() {
	s.moves = make(path.Moves, 0)
	s.currentMove = 0
	s.prevScore = 0
}

func (s *PathFollowingSolver) Ticks() int {
	return 0
}

func (s *PathFollowingSolver) NextMove(st *state.State) direction.Direction {
	regenPath := false
	switch s.regenPath {
	case RegenEveryTick:
		regenPath = true
		break
	case RegenNever:
		regenPath = len(s.moves) == 0
		break
	case RegenEveryFruit:
		regenPath = len(s.moves) == 0 || s.prevScore != st.Score()
		break
	}

	if regenPath {
		targetPath, foundPath := s.pathGen.Generate(st, st.SnakeHead(), st.Fruit())
		if !foundPath {
			return direction.None
		}
		s.moves = path.PathToMoves(targetPath)
		s.currentMove = 0
		s.prevScore = st.Score()
	}

	nextMove := s.moves[s.currentMove]

	s.currentMove++
	if s.currentMove > len(s.moves)-1 {
		s.currentMove = 0
	}
	return nextMove
}
