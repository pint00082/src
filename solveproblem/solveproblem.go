package solveproblem

import (
	"caching"
	"fmt"
	"math"
	"math/rand"
	"models"
	"runtime"
	"savesolution"
	"sort"
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
	startingSolution := models.Solution{}
	startingSolution.Solution = newSolution(load, weight)
	startingSolution.Score = caching.CalculateCostWCache(startingSolution.Solution, load)
	for k := range startingSolution.Solution {
		sort.Ints(startingSolution.Solution[k])
	}
	simulatedAnnealing(load, weight, startingSolution, finishTime)
}

func simulatedAnnealing(load *models.Loaded, weight []models.Weight, currentSolution models.Solution, finishTime time.Time) {

	maxProcs := runtime.GOMAXPROCS(0) + runtime.GOMAXPROCS(0)/4
	//Create Random Pools
	//Main thread pool
	randomness := rand.New(rand.NewSource(time.Now().UnixNano()))

	//Goroutines pools
	randomWorkers := []*rand.Rand{}
	for i := 0; i < maxProcs; i++ {
		randomWorkers = append(randomWorkers, rand.New(rand.NewSource(time.Now().UnixNano()*int64(i+3))))
	}

	maxThermalPeriod := load.Periods
	//divide time into periods
	remainingTime := finishTime.Sub(time.Now())
	raiseAfter := (int64(remainingTime) / int64(load.Periods))
	thermalPeriod := 1
	timeToRaise := time.Now().Add(time.Duration(raiseAfter))
	bestOverallSolution := currentSolution
	for {
		//Generate Candidates
		var wg sync.WaitGroup
		candidateSolutions := make([]models.Solution, maxProcs)
		remainingThermalPeriods := maxThermalPeriod - thermalPeriod
		for i := 0; i < maxProcs; i++ {
			wg.Add(1)
			cs := currentSolution
			go generateCandidates(cs, remainingThermalPeriods, randomWorkers[i], candidateSolutions, i, load, &wg)
		}
		wg.Wait()
		//Find Best Candidate
		bestCandidate := 0
		bestCandidateScore := candidateSolutions[0].Score
		for i := 1; i < maxProcs; i++ {
			candidateScore := candidateSolutions[i].Score
			if candidateScore < bestCandidateScore {
				bestCandidate = i
				bestCandidateScore = candidateScore

			}
		}

		//Allow best Candidate to become current solution

		if bestCandidateScore < currentSolution.Score {
			currentSolution = candidateSolutions[bestCandidate]
			if currentSolution.Score < bestOverallSolution.Score {
				bestOverallSolution = currentSolution
			}
		} else {

			scoreDifference := float64(bestCandidateScore-currentSolution.Score) / float64(currentSolution.Score)
			exp := math.Exp(-scoreDifference)
			ln := -math.Log(float64(thermalPeriod) / float64(maxThermalPeriod))
			if exp*ln >= randomness.Float64() {
				currentSolution = candidateSolutions[bestCandidate]
			}

		}

		//raise thermal perriod
		if time.Now().After(timeToRaise) {
			timeToRaise = time.Now().Add(time.Duration(raiseAfter))
			if thermalPeriod == maxThermalPeriod {
				break
			} else {
				if currentSolution.Score > bestOverallSolution.Score {
					currentSolution = bestOverallSolution
				}
				thermalPeriod++
				fmt.Println("entering stage:", thermalPeriod, " of ", maxThermalPeriod, " score: ", currentSolution.Score)
			}
		}
	}

	if currentSolution.Score <= bestOverallSolution.Score {
		currentSolution = bestOverallSolution
	}
	//report and save solution
	final := float64(currentSolution.Score) / float64(load.Students)
	final = math.Round(final*100) / 100
	savesolution.SaveSolution(load, currentSolution, final)
}
