package main

import (
	"flag"
	"fmt"
	"github.com/jmwri/neatgo/neat"
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

	if *slvr == "neat" {
		runNeat(writer)
		return
	}

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

func runNeat(writer score.Writer) {
	cfg := neat.DefaultConfig(12, 4)
	cfg.PopulationSize = 150

	cfg.BiasNodes = 1

	cfg.AddNodeMutationRate = .5
	cfg.DeleteNodeMutationRate = .3
	cfg.BiasMutationRate = .5
	cfg.BiasMutationPower = .3
	cfg.BiasReplaceRate = .3

	cfg.AddConnectionMutationRate = .7
	cfg.DeleteConnectionMutationRate = .4
	cfg.WeightMutationRate = .5
	cfg.WeightMutationPower = .3
	cfg.WeightReplaceRate = .3

	cfg.SpeciesCompatExcessCoeff = 1
	cfg.SpeciesCompatBiasDiffCoeff = .5
	cfg.SpeciesCompatWeightDiffCoeff = .5
	cfg.SpeciesCompatThreshold = 4.5
	cfg.SpeciesStalenessThreshold = 15
	cfg.MateCrossoverRate = .6
	cfg.MateBestRate = .5

	cfg.TopGenomesFromSpeciesToFill = 2

	pop, err := neat.GeneratePopulation(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting batch")
	for i := 0; i < *batchSize; i++ {
		fmt.Println("starting generation", i)
		clientStates := pop.States()
		wg := sync.WaitGroup{}
		wg.Add(len(clientStates))
		for _, state := range clientStates {
			go func(state neat.ClientGenomeState) {
				defer wg.Done()
				gameSolver := solver.NewNeatSolver(state)
				gameRunner := runner.NewGameRunner(*gameSize, *gameSize, gameSolver, writer)
				gameRunner.Init()

				movesNotImproved := 0
				oldScore := 0
				for {
					err := gameRunner.Update()
					if err != nil {
						fmt.Printf("%d failed: %s\n", i, err)
					}

					if gameRunner.State().Score() > oldScore {
						oldScore = gameRunner.State().Score()
						movesNotImproved = 0
					} else {
						movesNotImproved++
					}
					finished := gameRunner.Ended()
					aliveTooLong := false
					maxMovesWithoutImprovement := 300
					if movesNotImproved >= maxMovesWithoutImprovement {
						// Kill it!
						finished = true
						aliveTooLong = true
					}

					if finished {
						// Stop sending inputs to NN
						close(state.SendInput())
						// Send the fitness
						fastEatAdjustment := 1.0 - (float64(movesNotImproved) / float64(maxMovesWithoutImprovement))
						if aliveTooLong {
							fastEatAdjustment = 0
						}
						fitness := float64(gameRunner.State().Score()) + fastEatAdjustment
						state.SendFitness() <- fitness
						close(state.SendFitness())
						return
					}
				}
			}(state)
		}
		pop = neat.RunGeneration(pop)
		wg.Wait()

		bestEverFitness := pop.BestEverGenomeFitness
		totFitness := .0
		for _, fitness := range pop.GenomeFitness {
			totFitness += fitness
		}
		avgFitness := totFitness / float64(len(pop.GenomeFitness))
		bestEverNumNodes := pop.BestEverGenome.NumNodes()
		bestEverNumConnections := pop.BestEverGenome.NumConnections()
		bestFitness := pop.BestGenomeFitness
		bestNumNodes := pop.BestGenome.NumNodes()
		bestNumConnections := pop.BestGenome.NumConnections()
		popSize := len(pop.Genomes)
		numSpecies := len(pop.Species)

		fmt.Printf(`Generation %d
AvgFitness: %f
BestEverFitness: %f
BestEverNodes: %d
BestEverConnections: %d
BestFitness: %f
BestNodes: %d
BestConnections: %d
PopSize: %d
NumSpecies: %d
-------------------------
`, i, avgFitness, bestEverFitness, bestEverNumNodes, bestEverNumConnections, bestFitness, bestNumNodes, bestNumConnections, popSize, numSpecies)
	}
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
