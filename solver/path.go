package solver

import (
	"go-snake-ai/direction"
	"go-snake-ai/path"
	"go-snake-ai/state"
)

func NewPathFollowingSolver(name string, pathGen path.Generator, genPathEveryTick bool) *PathFollowingSolver {
	return &PathFollowingSolver{
		name:             name,
		pathGen:          pathGen,
		genPathEveryTick: genPathEveryTick,
	}
}

type PathFollowingSolver struct {
	name             string
	pathGen          path.Generator
	genPathEveryTick bool
	moves            path.Moves
	currentMove      int
	prevScore        int
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
	return 5
}

func (s *PathFollowingSolver) NextMove(st *state.State) direction.Direction {
	regenPath := false
	if s.genPathEveryTick || len(s.moves) == 0 || s.prevScore != st.Score() {
		regenPath = true
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
