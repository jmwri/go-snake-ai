package main

import (
	"flag"
	"go-snake-ai/game"
	"go-snake-ai/path"
	"go-snake-ai/scene"
	"go-snake-ai/score"
	"go-snake-ai/solver"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten"
)

var slvr = flag.String("solver", "user", "which solver to use")
var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
var gameSize = flag.Int("size", 10, "size of the play area")

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

	var gameSolver solver.Solver
	if *slvr == "user" {
		gameSolver = solver.NewUserSolver()
	} else if *slvr == "shortest" {
		pathGen := path.NewBreadthFirstSearch()
		gameSolver = solver.NewPathFollowingSolver(*slvr, pathGen, solver.RegenEveryTick)
	} else if *slvr == "longest" {
		bfs := path.NewBreadthFirstSearch()
		pathGen := path.NewBreadthFirstSearchLongest(bfs)
		gameSolver = solver.NewPathFollowingSolver(*slvr, pathGen, solver.RegenEveryFruit)
	} else {
		panic("no solver found")
	}

	writer := score.NewCSV("scores")

	titleScene := scene.NewTitleScene()
	gameScene := scene.NewGameScene(*gameSize, *gameSize, gameSolver, writer)
	manager := scene.NewManager(500, 500, titleScene, gameScene)
	opts := game.Options{
		NumTilesX: *gameSize,
		NumTilesY: *gameSize,
		Manager:   manager,
	}
	g := game.NewGame(opts)

	ebiten.SetWindowSize(manager.ScreenWidth(), manager.ScreenHeight())
	ebiten.SetWindowTitle("Snake AI")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
