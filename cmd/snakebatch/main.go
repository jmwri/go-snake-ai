package main

import (
	"flag"
	"fmt"
	"go-snake-ai/path"
	"go-snake-ai/runner"
	"go-snake-ai/score"
	"go-snake-ai/solver"
	"log"
	"os"
	"runtime/pprof"
	"sync"
)

var slvr = flag.String("solver", "user", "which solver to use")
var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
var gameSize = flag.Int("size", 10, "size of the play area")
var batchSize = flag.Int("batchSize", 1000, "how many games to run in a batch")

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

	writer := score.NewCSV("scores")

	wg := sync.WaitGroup{}
	fmt.Println("starting batch")
	for i := 0; i < *batchSize; i++ {
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
		} else if *slvr == "hamiltonian" {
			bfs := path.NewBreadthFirstSearch()
			longest := path.NewBreadthFirstSearchLongest(bfs)
			pathGen := path.NewHamiltonianCycle(longest)
			gameSolver = solver.NewPathFollowingSolver(*slvr, pathGen, solver.RegenNever)
		} else {
			panic("no solver found")
		}

		wg.Add(1)
		go runGameInstance(i, &wg, writer, gameSolver)
	}

	wg.Wait()
	fmt.Println("finished batch")
}

func runGameInstance(i int, wg *sync.WaitGroup, writer score.Writer, gameSolver solver.Solver) {
	defer wg.Done()
	gameRunner := runner.NewGameRunner(*gameSize, *gameSize, gameSolver, writer)
	gameRunner.Init()
	for {
		err := gameRunner.Update()
		if err != nil {
			fmt.Printf("%d failed: %s\n", i, err)
		}

		if gameRunner.Ended() {
			fmt.Printf("%d finished with score %d\n", i, gameRunner.State().Score())
			return
		}
	}
}
