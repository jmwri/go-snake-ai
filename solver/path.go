package solver

import (
	"go-snake-ai/direction"
	"go-snake-ai/path"
	"go-snake-ai/state"
	"go-snake-ai/tile"
	"math/rand"
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
			return s.randomDirection(st)
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

func (s *PathFollowingSolver) randomDirection(st *state.State) direction.Direction {
	// If we couldn't find a path, move in a random direction that won't kill the snake
	adjVectors := tile.AdjacentVectors(st.SnakeHead())
	freeDirections := make([]direction.Direction, 0)
	for dir, adjV := range adjVectors {
		if !st.ValidPosition(adjV.X, adjV.Y) {
			continue
		}
		t := st.Tile(adjV.X, adjV.Y)
		if t == tile.TypeBody || t == tile.TypeHead {
			continue
		}
		if direction.IsOpposite(dir, st.SnakeDir()) {
			continue
		}
		freeDirections = append(freeDirections, dir)
	}
	if len(freeDirections) == 0 {
		return direction.None
	}
	randomI := rand.Intn(len(freeDirections))
	return freeDirections[randomI]
}
