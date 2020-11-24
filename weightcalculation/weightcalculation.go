package weightcalculation

import (
	"models"
	"sort"
	"sync"
)

//CalculateWeights finds the Weigth of each Lesson
func CalculateWeights(load *models.Loaded) []models.Weight {
	lessons := load.Lessons
	weight := make([]int, lessons)
	for i := 0; i < lessons; i++ {
		weight[i] = 0
	}

	w := 1
calculate:
	for {
		degree := make([]int, lessons)
		for i := 0; i < lessons; i++ {
			degree[i] = 0
		}

		var wg sync.WaitGroup
		for i := 1; i <= lessons; i++ {
			wg.Add(1)
			go findDegree(&weight, &degree, load, i, w, &wg)
		}
		wg.Wait()

		minimumDegree := 5000
		for _, v := range degree {
			if v < minimumDegree {
				minimumDegree = v
			}
		}
		change := false
		for key, v := range degree {
			if v == minimumDegree || v == 0 {
				weight[key] = w
				change = true
			}
		}
		if change == true {
			w++
		}
		ready := true
		for _, v := range weight {
			if v == 0 {
				ready = false
			}
		}
		if ready == true {
			break calculate
		}
	}
	weighDist := distributeWeights(weight)
	return weighDist
}

func findDegree(weight, degree *[]int, load *models.Loaded, i int, w int, wg *sync.WaitGroup) {
	defer wg.Done()

	for k, common := range load.GraphNodes[i] {
		if (*weight)[i-1] != 0 {
			(*degree)[i-1] = 5000
		}
		if (*weight)[i-1] == 0 && common != 0 && (*weight)[k-1] <= w {
			(*degree)[i-1]++
		}
	}
}

func distributeWeights(weight []int) []models.Weight {
	distWeight := make([]models.Weight, len(weight))
	for k, v := range weight {
		m := models.Weight{}
		m.Lesson = k + 1
		m.Weight = v
		distWeight = append(distWeight, m)
	}
	sort.Slice(distWeight, func(i, j int) bool { return distWeight[i].Weight > distWeight[j].Weight })

	return distWeight
}
