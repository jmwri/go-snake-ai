package state

import (
	"fmt"
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
		totTiles:     tileNumX * tileNumY,
		tiles:        tiles,
		changedTiles: changed,
		score:        0,
		maxScore:     (tileNumX * tileNumY) - 1,
	}
	s.spawnSnake()
	s.spawnFruit()
	return s
}

// State represents a game field with tile states.
type State struct {
	tileNumX     int
	tileNumY     int
	totTiles     int
	tiles        [][]tile.Type
	changedTiles []*tile.Vector
	snakeDir     direction.Direction
	snake        []*tile.Vector
	fruit        *tile.Vector
	score        int
	maxScore     int
}

func (s *State) ValidPosition(x int, y int) bool {
	if x < 0 || x > s.tileNumX-1 || y < 0 || y > s.tileNumY-1 {
		return false
	}
	return true
}

// SetTile sets the tile at the given coordinates
func (s *State) SetTile(x int, y int, tileType tile.Type) {
	s.tiles[y][x] = tileType
}

// Tile returns the tile from the given vector
func (s *State) Tile(x int, y int) tile.Type {
	return s.tiles[y][x]
}

// Tiles returns the tiles from the state
func (s *State) Tiles() [][]tile.Type {
	return s.tiles
}

func (s *State) Score() int {
	return s.score
}

func (s *State) TotalTiles() int {
	return s.totTiles
}

func (s *State) MaxScore() int {
	return s.maxScore
}

func (s *State) Won() bool {
	return s.score == s.maxScore
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

func (s *State) spawnSnake() {
	snakePos := s.randomFreeTile()
	s.snake = []*tile.Vector{snakePos}
	s.SetTile(snakePos.X, snakePos.Y, tile.TypeHead)

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

func (s *State) spawnFruit() bool {
	fruitPos := s.randomFreeTile()
	if fruitPos == nil {
		return false
	}
	s.SetTile(fruitPos.X, fruitPos.Y, tile.TypeFruit)
	s.fruit = fruitPos
	return true
}

func (s *State) Move(dir direction.Direction) (bool, error) {
	curVec := s.SnakeHead()
	var nextVec *tile.Vector

	if dir == direction.None || direction.IsOpposite(dir, s.snakeDir) {
		dir = s.snakeDir
	}

	s.snakeDir = dir

	if dir == direction.Up {
		nextVec = tile.NewVector(curVec.X, curVec.Y-1)
	} else if dir == direction.Right {
		nextVec = tile.NewVector(curVec.X+1, curVec.Y)
	} else if dir == direction.Down {
		nextVec = tile.NewVector(curVec.X, curVec.Y+1)
	} else if dir == direction.Left {
		nextVec = tile.NewVector(curVec.X-1, curVec.Y)
	} else {
		return false, fmt.Errorf("unable to move")
	}

	if nextVec.X < 0 || nextVec.X >= s.tileNumX || nextVec.Y < 0 || nextVec.Y >= s.tileNumY {
		return false, nil
	}

	targetTile := s.Tile(nextVec.X, nextVec.Y)
	if targetTile == tile.TypeFruit {
		s.extendSnake(nextVec, false)
		// Increment score
		s.score++
		// Respawn fruit
		if !s.spawnFruit() {
			return false, nil
		}

		return true, nil
	} else if targetTile == tile.TypeNone {
		s.extendSnake(nextVec, true)
		return true, nil
	}
	return false, nil
}

func (s *State) extendSnake(next *tile.Vector, removeTail bool) {
	// Keep track of old head as we need to change it to TypeBody
	oldHead := s.SnakeHead()
	s.SetTile(oldHead.X, oldHead.Y, tile.TypeBody)
	// Eating a fruit! Set the fruit tile to snake body
	s.SetTile(next.X, next.Y, tile.TypeHead)
	// Set new head of snake to the next vector
	s.snake = append([]*tile.Vector{next}, s.snake...)

	if removeTail {
		// Remove tail of snake
		tailVec := s.SnakeTail()
		// Set tile of tail to None
		s.SetTile(tailVec.X, tailVec.Y, tile.TypeNone)
		// Remove last vector from snake slice
		s.snake = s.snake[:len(s.snake)-1]
	}
}

func (s *State) Fruit() *tile.Vector {
	return s.fruit
}

func (s *State) SnakeDir() direction.Direction {
	return s.snakeDir
}

func (s *State) SnakeHead() *tile.Vector {
	return s.snake[0]
}

func (s *State) SnakeTail() *tile.Vector {
	snakeLen := len(s.snake)
	return s.snake[snakeLen-1]
}
