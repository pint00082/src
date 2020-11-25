package solveproblem

import (
	"climber"
	"fmt"
	"math"
	"models"
	"newclimber"
	"runtime"
	"savesolution"
	"sync"
	"time"
	weight "weightcalculation"
)

//SolveProblem searches for the best Solution
func SolveProblem(load *models.Loaded, seconds int) {
	//clean systems cache
	load.Cache.Purge()
	//find when to finish
	finishTime := time.Now().Add(time.Duration(seconds) * time.Second)
	//distribute weights
	weight := weight.CalculateWeights(load)

	//greedyColoring
	currentSolution := models.Solution{}
	currentSolution.Solution = newclimber.NewClimber(load, weight)
	currentSolution.Score = climber.FindCost(currentSolution.Solution, load)

	//begin HillClimbing
	climbHill(load, weight, currentSolution, finishTime)
}

func climbHill(load *models.Loaded, weight []models.Weight, currentSolution models.Solution, finishTime time.Time) {
	maxProcs := runtime.GOMAXPROCS(0) + runtime.GOMAXPROCS(0)/4
	var wg sync.WaitGroup
	workerSolutionChan := make(chan models.Solution, maxProcs)
	defer close(workerSolutionChan)
	var steps int
	steps = load.Conf.StepFactor
	var raiseStepAfter int
	raiseStepAfter = load.Conf.StepFactor
	unchangedCounter := 0
	stuck := 0
	contestants := []models.Solution{}

	//main Loop
	for time.Now().Before(finishTime) {

		workerResults := []models.Solution{}
		//start the workers
		for i := 0; i < maxProcs; i++ {
			wg.Add(1)
			go climber.Climbing(load, steps, currentSolution, i, workerSolutionChan, &wg)

		}
		//wait for the worker results
		for i := 0; i < maxProcs; i++ {
			workerResults = append(workerResults, <-workerSolutionChan)

		}
		//ensure that all workers finished
		wg.Wait()
		// check for the best solution the workers found
		var minPos int
		minPos = 0
		bestSolution := workerResults[0].Score
		for k, v := range workerResults {
			if v.Score < bestSolution {
				minPos = k
				bestSolution = v.Score
			}
		}
		//if workers found a better solution keep it as the current best, resets the step factor and unchangedCounter
		if bestSolution < currentSolution.Score {
			currentSolution = workerResults[minPos]
			//steps = load.Conf.StepFactor
			unchangedCounter = 0
			//fmt.Println(currentSolution.Score)
		} else {
			unchangedCounter++
		}
		//raise  stepFactor if current solution hasn't changed since raiseStep times
		if unchangedCounter >= raiseStepAfter && steps < load.Conf.MaxSteps {
			steps++
			unchangedCounter = 0
		}
		if steps >= load.Conf.MaxSteps && unchangedCounter >= raiseStepAfter {
			stuck++
			//keepCurrentSolution and restart if maximum steps are reached and there has been no change for conf.Restart times
			if stuck == load.Conf.Restart {
				contestants = append(contestants, currentSolution)
				steps = load.Conf.StepFactor
				unchangedCounter = 0
				stuck = 0
				currentSolution.Solution = newclimber.NewClimber(load, weight)
				currentSolution.Score = climber.FindCost(currentSolution.Solution, load)
			}
		}
	}

	//find the best solution in the final contestants
	if len(contestants) > 0 {
		var minPos int
		minPos = 0
		bestSolution := contestants[0].Score
		for k, v := range contestants {
			if v.Score < bestSolution {
				minPos = k
				bestSolution = v.Score
			}
		}
		currentSolution = contestants[minPos]
	}
	//report and save solution
	//TODO
	final := float64(currentSolution.Score) / float64(load.Students)
	final = math.Round(final*100) / 100
	fmt.Println("score: ", final)
	savesolution.SaveSolution(load, currentSolution, final)
}
