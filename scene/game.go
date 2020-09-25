package scene

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"go-snake-ai/state"
	"go-snake-ai/ui"
	"math/rand"
	"time"
)

func NewGameScene(tileNumX int, tileNumY int) *GameScene {
	return &GameScene{
		manager:       nil,
		s:             nil,
		score:         0,
		gameover:      false,
		imageScore:    nil,
		imageGameover: nil,
		tileNumX:      tileNumX,
		tileNumY:      tileNumY,
		tileWidth:     0,
		tileHeight:    0,
	}
}

type GameScene struct {
	manager       *Manager
	s             *state.State
	score         int
	gameover      bool
	imageScore    *ebiten.Image
	imageGameover *ebiten.Image
	tileNumX      int
	tileNumY      int
	tileWidth     float64
	tileHeight    float64
}

func (s *GameScene) SetManager(manager *Manager) {
	s.manager = manager
}

func (s *GameScene) Init() {
	rand.Seed(time.Now().UnixNano())
	s.imageScore, _ = ebiten.NewImage(s.manager.ScreenWidth(), s.manager.ScreenHeight(), ebiten.FilterDefault)
	s.imageGameover, _ = ebiten.NewImage(s.manager.ScreenWidth(), s.manager.ScreenHeight(), ebiten.FilterDefault)
	s.tileWidth = float64(s.manager.ScreenWidth()) / float64(s.tileNumX)
	s.tileHeight = float64(s.manager.ScreenHeight()) / float64(s.tileNumY)
	s.s = state.NewState(s.tileNumX, s.tileNumY)
}

func (s *GameScene) Update() error {
	if s.gameover {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			s.manager.GoToTitle()
		}
		return nil
	}
	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	for y, col := range s.s.Tiles() {
		for x, t := range col {
			ui.DrawRect(screen, s.tileWidth*float64(x), s.tileHeight*float64(y), s.tileWidth, s.tileHeight, ui.TileColours[t])
		}
	}
}
