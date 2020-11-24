package climber

import (
	"caching"
	"checksolution"
	"math/rand"
	"models"
	"move"
	"sync"
	"time"
)

//Climbing is the main process
func Climbing(load *models.Loaded, stepFactor int, currentSolution models.Solution, w int, workerSolutionChan chan models.Solution, wg *sync.WaitGroup) {
	defer wg.Done()
	random := rand.New(rand.NewSource(time.Now().UnixNano() * int64(w)))
	//create random sequence of moves

	climber := make([][]int, len(currentSolution.Solution))
	for i := range currentSolution.Solution {
		climber[i] = make([]int, len(currentSolution.Solution[i]))
		copy(climber[i], currentSolution.Solution[i])
	}
	for i := 0; i < stepFactor; i++ {
		selectMove := random.Intn(load.MovesDist.PossibilitySum)
		if selectMove < load.MovesDist.ChangePeriods {
			move.ChangePeriods(load, climber, random)
		} else if selectMove < load.MovesDist.MoveExam {
			move.ExamMove(load, climber, random)
		} else if selectMove < load.MovesDist.ExchangeExams {
			move.ExchangeExams(load, climber, random)
		} else if selectMove < load.MovesDist.ExchangeExams {
			move.MassChangePeriods(load, climber, random)
		} else if selectMove < load.MovesDist.MassExodus {
			move.MassExodus(load, climber, random)
		} else if selectMove < load.MovesDist.MassExodus {
			//move.MassMigration(load, climber, random)
		}
	}
	workerSolution := models.Solution{}
	workerSolution.Solution = climber
	workerSolution.Score = CalculateCostWCache(climber, load)
	workerSolutionChan <- workerSolution
}

//FindCost to find total cost
func FindCost(solution [][]int, load *models.Loaded) int {
	periods := make(map[int][]int)
	for k, v := range solution {
		for _, value := range v {
			periods[k] = append(periods[k], int(value))
		}
	}
	return checksolution.CalculateCost(periods, load)
}

//CalculateCostWCache returns the Cost using the ArcCache
func CalculateCostWCache(solution [][]int, load *models.Loaded) int {
	totalCost := 0
	for i := 1; i < load.Periods; i++ {
		maxImpact := i - 5
		if maxImpact < 0 {
			maxImpact = 0
		}
		for j := i - 1; j >= maxImpact; j-- {
			periodCost := caching.PeriodClashes(solution[i], solution[j], load)
			totalCost += periodCost * checksolution.FindCostFactor(i, j)
		}

	}
	return totalCost
}
