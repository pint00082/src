package caching

import (
	"checksolution"
	"fmt"
	"models"
	"strings"
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

//CheckAlreadyTriedWeights to Add a weightDistribution or return false if this weightDistribution already exists
func CheckAlreadyTriedWeights(layerShuffleDistWeight []int, load *models.Loaded) bool {
	//turn the weightDistribution to string
	layerShuffleDistWeightString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(layerShuffleDistWeight)), ","), "[]")
	//if it already exists in the cache return false
	_, ok := load.Cache.Get(layerShuffleDistWeightString)
	if ok {
		return false
	}
	//else add to cache and return true
	load.Cache.Add(layerShuffleDistWeightString, "check")
	return true
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
			periodCost := PeriodClashes(solution[i], solution[j], load)
			totalCost += periodCost * checksolution.FindCostFactor(i, j)
		}

	}
	return totalCost
}
