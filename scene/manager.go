package scene

import (
	"github.com/hajimehoshi/ebiten"
)

func NewManager(screenWidth int, screenHeight int, titleScene Scene, gameScene Scene) *Manager {
	m := &Manager{
		screenWidth:        screenWidth,
		screenHeight:       screenHeight,
		current:            nil,
		next:               nil,
		transitionCount:    0,
		maxTransitionCount: 20,
		transitionFrom:     nil,
		transitionTo:       nil,
		titleScene:         titleScene,
		gameScene:          gameScene,
	}
	m.init()
	return m
}

type Manager struct {
	screenWidth        int
	screenHeight       int
	current            Scene
	next               Scene
	transitionCount    int
	maxTransitionCount int
	transitionFrom     *ebiten.Image
	transitionTo       *ebiten.Image
	titleScene         Scene
	gameScene          Scene
}

func (m *Manager) init() {
	m.titleScene.SetManager(m)
	m.titleScene.Init()
	m.gameScene.SetManager(m)
	m.gameScene.Init()
	m.transitionFrom, _ = ebiten.NewImage(m.screenWidth, m.screenHeight, ebiten.FilterDefault)
	m.transitionTo, _ = ebiten.NewImage(m.screenWidth, m.screenHeight, ebiten.FilterDefault)
}

func (m *Manager) ScreenWidth() int {
	return m.screenWidth
}

func (m *Manager) ScreenHeight() int {
	return m.screenHeight
}

func (m *Manager) Update() error {
	if m.transitionCount == 0 {
		return m.current.Update()
	}

	m.transitionCount--
	if m.transitionCount > 0 {
		return nil
	}

	m.current = m.next
	m.next = nil
	m.current.Init()
	return nil
}

func (m *Manager) Draw(r *ebiten.Image) {
	if m.transitionCount == 0 {
		m.current.Draw(r)
		return
	}

	m.transitionFrom.Clear()
	m.current.Draw(m.transitionFrom)

	m.transitionTo.Clear()
	m.next.Draw(m.transitionTo)

	r.DrawImage(m.transitionFrom, nil)

	alpha := 1 - float64(m.transitionCount)/float64(m.maxTransitionCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, alpha)
	r.DrawImage(m.transitionTo, op)
}

func (m *Manager) goTo(scene Scene) {
	if m.current == nil {
		m.current = scene
	} else {
		m.next = scene
		m.transitionCount = m.maxTransitionCount
	}
}

func (m *Manager) GoToTitle() {
	m.goTo(m.titleScene)
}

func (m *Manager) GoToGame() {
	m.goTo(m.gameScene)
}
