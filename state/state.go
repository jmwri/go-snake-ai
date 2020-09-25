package state

import "go-snake-ai/tile"

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
	return &State{
		tileNumX:     tileNumX,
		tileNumY:     tileNumY,
		tiles:        tiles,
		changedTiles: changed,
	}
}

// State represents a game field with tile states.
type State struct {
	tileNumX     int
	tileNumY     int
	tiles        [][]tile.Type
	changedTiles []*tile.Vector
}

// SetTile sets the tile at the given coordinates
func (s *State) SetTile(x int, y int, tileType tile.Type) {
	s.tiles[x][y] = tileType
}

// Tiles returns the tiles from the state
func (s *State) Tiles() [][]tile.Type {
	return s.tiles
}
