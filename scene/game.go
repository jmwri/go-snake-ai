package scene

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"go-snake-ai/runner"
	"go-snake-ai/ui"
	"image/color"
	"strconv"
	"time"
)

func NewGameScene(runner *runner.GameRunner, updateTicks int) *GameScene {
	return &GameScene{
		runner:      runner,
		updateTicks: updateTicks,
		listening:   false,
	}
}

type GameScene struct {
	runner      *runner.GameRunner
	manager     *Manager
	imageScore  *ebiten.Image
	imageEnded  *ebiten.Image
	tileWidth   float64
	tileHeight  float64
	currentTick int
	updateTicks int
	listening   bool
}

func (s *GameScene) SetManager(manager *Manager) {
	s.manager = manager
}

func (s *GameScene) Init() {
	s.runner.Init()
	s.imageScore, _ = ebiten.NewImage(s.manager.ScreenWidth(), s.manager.ScreenHeight(), ebiten.FilterDefault)
	s.imageEnded, _ = ebiten.NewImage(s.manager.ScreenWidth(), s.manager.ScreenHeight(), ebiten.FilterDefault)
	s.tileWidth = float64(s.manager.ScreenWidth()) / float64(s.runner.TileNumX())
	s.tileHeight = float64(s.manager.ScreenHeight()) / float64(s.runner.TileNumY())
	s.currentTick = 0
	if !s.listening {
		s.listening = true
		go s.listenSpeed()
	}
}

func (s *GameScene) listenSpeed() {
	for {
		if ebiten.IsKeyPressed(ebiten.KeyComma) {
			s.updateTicks++
		} else if ebiten.IsKeyPressed(ebiten.KeyPeriod) && s.updateTicks > 0 {
			s.updateTicks--
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (s *GameScene) Update() error {
	s.currentTick++
	if s.runner.Ended() {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			s.manager.GoToTitle()
		}
		return nil
	}

	if s.currentTick < s.updateTicks {
		return nil
	}
	err := s.runner.Update()
	if err != nil {
		return err
	}

	s.currentTick = 0

	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	if s.runner.Ended() {
		endText := "game over"
		if s.runner.State().Won() {
			endText = "you won!"
		}
		ui.DrawTextWithShadowCenter(screen, endText, 0, 32, 4, color.NRGBA{0x00, 0x00, 0x80, 0xff}, s.manager.ScreenWidth())
		scoreText := fmt.Sprintf("score: %d", s.runner.State().Score())
		ui.DrawTextWithShadowCenter(screen, scoreText, 0, 100, 4, color.NRGBA{0x00, 0x00, 0x80, 0xff}, s.manager.ScreenWidth())
	} else {
		for y, col := range s.runner.State().Tiles() {
			for x, t := range col {
				ui.DrawRect(screen, s.tileWidth*float64(x), s.tileHeight*float64(y), s.tileWidth, s.tileHeight, ui.TileColours[t])
			}
		}

		// Draw score
		ui.DrawTextWithShadow(screen, strconv.Itoa(s.runner.State().Score()), 10, 10, 1, color.White)
	}
}
