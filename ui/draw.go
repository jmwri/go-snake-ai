package ui

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
)

func DrawRect(dst *ebiten.Image, x, y, width, height float64, clr color.Color)  {
	ebitenutil.DrawRect(dst, x, y, width, height, clr)
}
