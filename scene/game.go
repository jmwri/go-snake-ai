package scene

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"go-snake-ai/score"
	"go-snake-ai/solver"
	"go-snake-ai/state"
	"go-snake-ai/ui"
	"image/color"
	"math/rand"
	"strconv"
	"time"
)

func NewGameScene(tileNumX int, tileNumY int, slvr solver.Solver, writer score.Writer) *GameScene {
	return &GameScene{
		tileNumX: tileNumX,
		tileNumY: tileNumY,
		slvr:     slvr,
		writer:   writer,
	}
}

type GameScene struct {
	manager      *Manager
	s            *state.State
	ended        bool
	imageScore   *ebiten.Image
	imageEnded   *ebiten.Image
	tileNumX     int
	tileNumY     int
	tileWidth    float64
	tileHeight   float64
	slvr         solver.Solver
	currentTick  int
	writer       score.Writer
	writtenScore bool
}

func (s *GameScene) SetManager(manager *Manager) {
	s.manager = manager
}

func (s *GameScene) Init() {
	rand.Seed(time.Now().UnixNano())
	s.imageScore, _ = ebiten.NewImage(s.manager.ScreenWidth(), s.manager.ScreenHeight(), ebiten.FilterDefault)
	s.imageEnded, _ = ebiten.NewImage(s.manager.ScreenWidth(), s.manager.ScreenHeight(), ebiten.FilterDefault)
	s.tileWidth = float64(s.manager.ScreenWidth()) / float64(s.tileNumX)
	s.tileHeight = float64(s.manager.ScreenHeight()) / float64(s.tileNumY)
	s.s = state.NewState(s.tileNumX, s.tileNumY)
	s.slvr.Init()
	s.ended = false
	s.currentTick = 0
	s.writtenScore = false
}

func (s *GameScene) Update() error {
	s.currentTick++
	if s.ended {
		if !s.writtenScore {
			s.writtenScore = true
			s.writer.Write(s.s.Score(), s.s.MaxScore(), s.slvr)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			s.manager.GoToTitle()
		}
		return nil
	}

	if s.currentTick < s.slvr.Ticks() {
		return nil
	}
	nextDirection := s.slvr.NextMove(s.s)
	alive, err := s.s.Move(nextDirection)
	if err != nil {
		return err
	}
	if !alive {
		s.ended = true
	}

	s.currentTick = 0

	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	if s.ended {
		endText := "game over"
		if s.s.Won() {
			endText = "you won!"
		}
		ui.DrawTextWithShadowCenter(screen, endText, 0, 32, 4, color.NRGBA{0x00, 0x00, 0x80, 0xff}, s.manager.ScreenWidth())
		scoreText := fmt.Sprintf("score: %d", s.s.Score())
		ui.DrawTextWithShadowCenter(screen, scoreText, 0, 100, 4, color.NRGBA{0x00, 0x00, 0x80, 0xff}, s.manager.ScreenWidth())
	} else {
		for y, col := range s.s.Tiles() {
			for x, t := range col {
				ui.DrawRect(screen, s.tileWidth*float64(x), s.tileHeight*float64(y), s.tileWidth, s.tileHeight, ui.TileColours[t])
			}
		}

		// Draw score
		ui.DrawTextWithShadow(screen, strconv.Itoa(s.s.Score()), 10, 10, 1, color.White)
	}
}