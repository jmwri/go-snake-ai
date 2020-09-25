package main

import (
	"flag"
	"go-snake-ai/game"
	"go-snake-ai/scene"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	titleScene := scene.NewTitleScene()
	gameScene := scene.NewGameScene(10, 10)
	manager := scene.NewManager(500, 500, titleScene, gameScene)
	opts := game.Options{
		NumTilesX: 10,
		NumTilesY: 10,
		Manager:   manager,
	}
	g := game.NewGame(opts)

	ebiten.SetWindowSize(manager.ScreenWidth(), manager.ScreenHeight())
	ebiten.SetWindowTitle("Snake AI")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
