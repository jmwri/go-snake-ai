package scene

import (
	"go-snake-ai/ui"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func NewTitleScene() *TitleScene {
	return &TitleScene{
		manager: nil,
	}
}

type TitleScene struct {
	manager *Manager
}

func (s *TitleScene) Init() {

}

func (s *TitleScene) SetManager(manager *Manager) {
	s.manager = manager
}

func (s *TitleScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.manager.GoToGame()
		return nil
	}

	return nil
}

func (s *TitleScene) Draw(r *ebiten.Image) {
	s.drawLogo(r, "SNAKE")

	message := "PRESS SPACE TO START"
	x := 0
	y := s.manager.ScreenHeight() - 48
	ui.DrawTextWithShadowCenter(r, message, x, y, 1, color.NRGBA{0x80, 0, 0, 0xff}, s.manager.ScreenWidth())
}

func (s *TitleScene) drawLogo(r *ebiten.Image, str string) {
	const scale = 4
	x := 0
	y := 32
	ui.DrawTextWithShadowCenter(r, str, x, y, scale, color.NRGBA{0x00, 0x00, 0x80, 0xff}, s.manager.ScreenWidth())
}
