package solveproblem

import (
	"caching"
	"fmt"
	"math/rand"
	"models"
	"runtime"
	"sync"
	"time"
)

//NewSolution finds a starting solution
func newSolution(load *models.Loaded, weight []models.Weight) [][]int {
	//Clean Cache
	load.Cache.Purge()
	colored := parallelColoring(load, weight)
	startingSolution := make([][]int, load.Periods)
	for k, v := range colored {
		startingSolution[v] = append(startingSolution[v], int(k+1))
	}
	return startingSolution
}

func parallelColoring(load *models.Loaded, weight []models.Weight) []int {
	var colored []int
	//find the number of threads in the system
	maxProcs := runtime.GOMAXPROCS(0)
	//create a WaitGroup to ensure that all threads finish
	var wg sync.WaitGroup
	//Create channels to force the workers to stop
	var suicideChans []chan bool
	//Create a channel to receive solutions
	solutionChannel := make(chan []int, 3*maxProcs)
	//start the workers
	for i := 0; i < maxProcs; i++ {
		wg.Add(1)
		suicideChans = append(suicideChans, make(chan bool))
		go parallelColoringWorkers(load, weight, i, suicideChans[i], solutionChannel, &wg)
	}

	//wait till first solution
	colored = <-solutionChannel
	fmt.Println("Found a starting Solution")
	//send signal for the workers to stop
	for i := 0; i < maxProcs; i++ {
		suicideChans[i] <- true
	}

	//ensure that all workers finished
	wg.Wait()
	return colored
}

func parallelColoringWorkers(load *models.Loaded, weight []models.Weight, w int, suicideChan chan bool, solutionChannel chan []int, wg *sync.WaitGroup) {
	//Sync the waitgroup
	defer wg.Done()
	//

	//Create a random pool
	random := rand.New(rand.NewSource(time.Now().UnixNano() * int64(w)))
	shuffle := load.Conf.Shuffle
	failedCounter := 0
	raiseShuffle := load.Periods * load.Conf.RaiseShuffleFactor
	maxShuffle := load.Conf.MaxShuffle
	//Check for termination signal
	//Check if there is a message in suicideChan to terminate the process
	for {

		select {
		case <-suicideChan:
			return
		default:
			if shuffle == maxShuffle {
				failedCounter = 0
			}
			if failedCounter >= raiseShuffle {
				shuffle++
				failedCounter = 0
			}
			//shuffle the weights
			layerShuffleDistWeight := layerShuffleWeight(shuffle, load, weight, random)
			//ensure that this layerShuffleDistWeight hasn't been tried in the past
			alreadyChecked := caching.CheckAlreadyTriedWeights(layerShuffleDistWeight, load)
			if alreadyChecked == false {
				failedCounter++
				continue
			}

			var colored []int
			ok := false
			colored, ok = greedyColoring(load, layerShuffleDistWeight)

			if ok == true {
				select {
				case solutionChannel <- colored:
				default:
				}
			}
			failedCounter++
		}
	}
}

func layerShuffleWeight(sh int, load *models.Loaded, weight []models.Weight, rouRand *rand.Rand) []int {

	randDistWeight := []models.Weight{}
	for _, v := range weight {
		randDistWeight = append(randDistWeight, v)
	}
	length := load.Lessons
	slice := []models.Weight{}
	shuffle := (length * sh) / 100
	if shuffle <= 0 {
		shuffle = 1
	}
	for i := 0; i < length; i += shuffle {
		var lastLessons int
		if i+shuffle >= length {
			lastLessons = length
		} else {
			lastLessons = i + shuffle
		}
		reslice := randDistWeight[i:lastLessons]
		rouRand.Shuffle(len(reslice), func(p1, p2 int) { reslice[p1], reslice[p2] = reslice[p2], reslice[p1] })
		for _, v := range reslice {
			slice = append(slice, v)
		}
	}
	seed := []int{}
	for _, v := range slice {
		seed = append(seed, v.Lesson)
	}
	return seed
}

func greedyColoring(load *models.Loaded, layerShuffleDistWeight []int) ([]int, bool) {
	length := load.Lessons
	colorLength := load.Periods
	colors := []int{}
	for i := 0; i < length; i++ {
		colors = append(colors, -1)
	}
	for i := 0; i < length; i++ {
		availableColors := []bool{}
		for c := 0; c < colorLength; c++ {
			availableColors = append(availableColors, true)
		}

		toBeColored := layerShuffleDistWeight[i]
		for k, v := range load.GraphNodes[toBeColored] {
			if v != 0 {
				if colors[k-1] != -1 {
					availableColors[colors[k-1]] = false
				}
			}
		}

		for c := 0; c < colorLength; c++ {
			if availableColors[c] == true {
				colors[toBeColored-1] = c
				break
			}
		}

		if colors[toBeColored-1] == -1 {
			return []int{}, false
		}
	}
	return colors, true
}
