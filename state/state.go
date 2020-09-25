package state

import (
	"go-snake-ai/direction"
	"go-snake-ai/tile"
	"math/rand"
)

// NewState returns a State that has been initialised with empty tiles
func NewState(tileNumX int, tileNumY int) *State {
	changed := make([]*tile.Vector, 0)
	tiles := make([][]tile.Type, tileNumY)
	for y := 0; y < tileNumY; y++ {
		row := make([]tile.Type, tileNumX)
		for x := 0; x < tileNumX; x++ {
			row[x] = tile.TypeNone
			changed = append(changed, tile.NewVector(x, y))
		}
		tiles[y] = row
	}
	s := &State{
		tileNumX:     tileNumX,
		tileNumY:     tileNumY,
		tiles:        tiles,
		changedTiles: changed,
	}
	s.SpawnSnake()
	s.SpawnFruit()
	return s
}

// State represents a game field with tile states.
type State struct {
	tileNumX     int
	tileNumY     int
	tiles        [][]tile.Type
	changedTiles []*tile.Vector
	snakeDir     direction.Direction
}

// SetTile sets the tile at the given coordinates
func (s *State) SetTile(x int, y int, tileType tile.Type) {
	s.tiles[x][y] = tileType
}

// Tile returns the tile from the given vector
func (s *State) Tile(x int, y int) tile.Type {
	return s.tiles[y][x]
}

// Tiles returns the tiles from the state
func (s *State) Tiles() [][]tile.Type {
	return s.tiles
}

func (s *State) randomFreeTile() *tile.Vector {
	freeVectors := make([]*tile.Vector, 0)
	for y, col := range s.tiles {
		for x, t := range col {
			if t == tile.TypeNone {
				freeVectors = append(freeVectors, tile.NewVector(x, y))
			}
		}
	}
	if len(freeVectors) == 0 {
		return nil
	}
	randomI := rand.Intn(len(freeVectors))
	return freeVectors[randomI]
}

func (s *State) SpawnSnake() {
	snakePos := s.randomFreeTile()
	s.SetTile(snakePos.X, snakePos.Y, tile.TypeBody)

	freeDirections := make([]direction.Direction, 0)
	if snakePos.X > 0 && s.Tile(snakePos.X-1, snakePos.Y) == tile.TypeNone {
		freeDirections = append(freeDirections, direction.Left)
	}
	if snakePos.X < s.tileNumX-1 && s.Tile(snakePos.X+1, snakePos.Y) == tile.TypeNone {
		freeDirections = append(freeDirections, direction.Right)
	}
	if snakePos.Y > 0 && s.Tile(snakePos.X, snakePos.Y-1) == tile.TypeNone {
		freeDirections = append(freeDirections, direction.Up)
	}
	if snakePos.Y < s.tileNumY-1 && s.Tile(snakePos.X, snakePos.Y+1) == tile.TypeNone {
		freeDirections = append(freeDirections, direction.Down)
	}

	randomI := rand.Intn(len(freeDirections))
	s.snakeDir = freeDirections[randomI]
}

func (s *State) SpawnFruit() {
	snakePos := s.randomFreeTile()
	s.SetTile(snakePos.X, snakePos.Y, tile.TypeFruit)
}
