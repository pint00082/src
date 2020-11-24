package newclimber

import (
	"math/rand"
	"models"
	"time"
)

//NewClimber shuffles the weight , finds a working starting Point
func NewClimber(load *models.Loaded, weight []models.Weight) [][]int {

	seed := rand.NewSource(time.Now().UnixNano())
	rouRand := rand.New(seed)
	randDistWeight := []models.Weight{}
	for _, v := range weight {
		randDistWeight = append(randDistWeight, v)
	}

	var colored []int
	ok := false
	for ok == false {
		layerShuffleDistWeight := layerShuffleWeight(load, randDistWeight, rouRand)
		colored, ok = greedyColoring(load, layerShuffleDistWeight)

	}
	climber := make([][]int, load.Periods)
	for k, v := range colored {
		climber[v] = append(climber[v], int(k+1))
	}
	rouRand.Shuffle(len(climber), func(i, j int) { climber[i], climber[j] = climber[j], climber[i] })
	return climber
}

func layerShuffleWeight(load *models.Loaded, randDistWeight []models.Weight, rouRand *rand.Rand) []int {
	lenght := load.Lessons
	slice := []models.Weight{}
	shuffle := (lenght * load.Conf.Shuffle) / 100
	if 20 >= lenght {
		shuffle = 1
	}
	for i := 0; i < lenght; i += shuffle {
		var lastLessons int
		if i+shuffle >= lenght {
			lastLessons = lenght
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
	colorLenght := load.Periods
	colors := []int{}
	for i := 0; i < length; i++ {
		colors = append(colors, -1)
	}
	for i := 0; i < length; i++ {
		availableColors := []bool{}
		for c := 0; c < colorLenght; c++ {
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

		for c := 0; c < colorLenght; c++ {
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
