package solveproblem

import (
	"caching"
	"math/rand"
	"models"
	"sync"
)

func generateCandidates(solution models.Solution, remainingThermalPeriods int, randPool *rand.Rand, candidateSolutions []models.Solution, w int, load *models.Loaded, wg *sync.WaitGroup) {
	defer wg.Done()

	currentSolution := models.Solution{}
	currentSolution.Solution = make([][]int, len(solution.Solution))
	for i := range solution.Solution {
		currentSolution.Solution[i] = make([]int, len(solution.Solution[i]))
		copy(currentSolution.Solution[i], solution.Solution[i])
	}
	currentSolution.Score = solution.Score
	maxPeriod := load.Periods
	//randomly choose how many moves will be executed
	moves := randPool.Intn(remainingThermalPeriods + 1)
	if moves == 0 {
		moves = 1
	}
	//Create the final form
	moveSequence := []int{}
	for i := 0; i < maxPeriod; i++ {
		moveSequence = append(moveSequence, i)
	}
	for i := 0; i < moves; i++ {
		periodA := 0
		periodB := 0
		//choose two different periods
		for periodA == periodB {
			periodA = randPool.Intn(maxPeriod)
			periodB = randPool.Intn(maxPeriod)
		}

		//roll for kempe
		kempeRoll := randPool.Intn(100) + 1
		if kempeRoll <= load.Conf.Kempe {
			kempe(currentSolution.Solution, load, randPool, moveSequence[periodA], moveSequence[periodB])
		} else {
			temp := moveSequence[periodA]
			moveSequence[periodA] = moveSequence[periodB]
			moveSequence[periodB] = temp
		}
	}
	//perform the exchange
	newCandidate := models.Solution{}
	newCandidate.Solution = [][]int{}
	for i := 0; i < maxPeriod; i++ {
		newCandidate.Solution = append(newCandidate.Solution, currentSolution.Solution[moveSequence[i]])
	}
	newCandidate.Score = caching.CalculateCostWCache(newCandidate.Solution, load)
	candidateSolutions[w] = newCandidate
}
