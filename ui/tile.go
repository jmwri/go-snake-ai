package ui

import (
	"go-snake-ai/tile"
	"image/color"
)

var TileColours = map[tile.Type]color.Color{
	tile.TypeNone: color.RGBA{
		R: 20,
		G: 20,
		B: 20,
		A: 255,
	},
	tile.TypeHead: color.RGBA{
		R: 0,
		G: 200,
		B: 0,
		A: 255,
	},
	tile.TypeBody: color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	},
	tile.TypeFruit: color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	},
}
