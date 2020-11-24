package caching

import (
	"fmt"
	"models"
)

//PeriodClashes returns the common students in two exams
func PeriodClashes(a []int, b []int, load *models.Loaded) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	var cacheK string
	if a[0] < b[0] {
		cacheK = fmt.Sprint(a) + " , " + fmt.Sprint(b)
	} else {
		cacheK = fmt.Sprint(b) + " , " + fmt.Sprint(a)
	}

	periodCost, ok := load.Cache.Get(cacheK)
	if ok {
		return periodCost.(int)
	}
	periodTotal := 0
	for _, v := range a {
		for _, value := range b {
			periodTotal += load.GraphNodes[int(v)][int(value)]
		}
	}
	load.Cache.Add(cacheK, periodTotal)
	return periodTotal
}
