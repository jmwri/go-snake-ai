package scene

import "github.com/hajimehoshi/ebiten"

type Scene interface {
	Init()
	SetManager(manager *Manager)
	Update() error
	Draw(screen *ebiten.Image)
}
