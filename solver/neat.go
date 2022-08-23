package solver

import (
	"fmt"
	"github.com/jmwri/neatgo/neat"
	"go-snake-ai/direction"
	"go-snake-ai/state"
	"go-snake-ai/tile"
	"math"
)

func NewNeatSolver(state neat.ClientGenomeState) *NeatSolver {
	return &NeatSolver{
		state: state,
	}
}

type NeatSolver struct {
	state neat.ClientGenomeState
}

func (s *NeatSolver) Name() string {
	return "neat"
}

func (s *NeatSolver) Init() {
}

func isTileType(t tile.Type, want []tile.Type) bool {
	for _, w := range want {
		if t == w {
			return true
		}
	}
	return false
}

func (s *NeatSolver) getProximity(st *state.State, want []tile.Type, x, y int) float64 {
	emptyBlocksInDir := .0
	currentPos := st.SnakeHead()
	for {
		currentPos = tile.NewVector(currentPos.X+x, currentPos.Y+y)
		if !st.ValidPosition(currentPos.X, currentPos.Y) {
			break
		}
		t := st.Tile(currentPos.X, currentPos.Y)
		if isTileType(t, want) {
			break
		}
		emptyBlocksInDir++
	}
	xTiles, yTiles := st.Dimensions()
	return emptyBlocksInDir / (float64(xTiles+yTiles) / 2)
}

func (s *NeatSolver) NextMove(st *state.State) direction.Direction {
	// Build inputs from state
	var badProximityUp, badProximityRight, badProximityDown, badProximityLeft float64
	badProximityUp = s.getProximity(st, []tile.Type{tile.TypeHead, tile.TypeBody}, 0, -1)
	badProximityRight = s.getProximity(st, []tile.Type{tile.TypeHead, tile.TypeBody}, 1, 0)
	badProximityDown = s.getProximity(st, []tile.Type{tile.TypeHead, tile.TypeBody}, 0, 1)
	badProximityLeft = s.getProximity(st, []tile.Type{tile.TypeHead, tile.TypeBody}, -1, 0)

	var goodProximityUp, goodProximityRight, goodProximityDown, goodProximityLeft float64
	goodProximityUp = s.getProximity(st, []tile.Type{tile.TypeFruit}, 0, -1)
	goodProximityRight = s.getProximity(st, []tile.Type{tile.TypeFruit}, 1, 0)
	goodProximityDown = s.getProximity(st, []tile.Type{tile.TypeFruit}, 0, 1)
	goodProximityLeft = s.getProximity(st, []tile.Type{tile.TypeFruit}, -1, 0)

	var headUp, headRight, headDown, headLeft float64
	switch st.SnakeDir() {
	case direction.Up:
		headUp = 1
		break
	case direction.Right:
		headRight = 1
		break
	case direction.Down:
		headDown = 1
		break
	case direction.Left:
		headLeft = 1
		break
	}
	inputs := []float64{
		badProximityUp,
		badProximityRight,
		badProximityDown,
		badProximityLeft,
		goodProximityUp,
		goodProximityRight,
		goodProximityDown,
		goodProximityLeft,
		headUp,
		headRight,
		headDown,
		headLeft,
	}

	// Send input to NN
	s.state.SendInput() <- inputs

	// Wait for outputs or error
	var outputs []float64
	select {
	case outputs = <-s.state.GetOutput():
		break
	case err := <-s.state.GetError():
		fmt.Printf("failed to process: %s\n", err)
		return direction.Up
	}

	// Return the best guessed direction
	var bestOutputIndex int
	bestOutput := math.Inf(-1)
	for i, output := range outputs {
		if output > bestOutput {
			bestOutputIndex = i
			bestOutput = output
		}
	}

	switch bestOutputIndex {
	case 0:
		return direction.Up
	case 1:
		return direction.Right
	case 2:
		return direction.Down
	case 3:
		return direction.Left
	}
	return direction.None
}
