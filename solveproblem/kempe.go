package solveproblem

import (
	"math/rand"
	"models"
	"sort"
)

//Kempe creates a kempe chain and executes the exchange
func kempe(solution [][]int, load *models.Loaded, random *rand.Rand, periodA int, periodB int) {
	perA := make([]int, len(solution[periodA]))
	copy(perA, solution[periodA])
	perB := make([]int, len(solution[periodB]))
	copy(perB, solution[periodB])

	//Check if periodA is empty
	if len(perA) == 0 {
		return
	}

	//Grab a random exam from period A
	exA := random.Intn(len(perA))
	examA := perA[exA]
	//create random sequence for periodB
	rand.Shuffle(len(perB), func(i, j int) { perB[i], perB[j] = perB[j], perB[i] })
	//find the first interconnected pair
	found := false
	var examB int
	for _, vB := range perB {
		if load.GraphNodes[examA][vB] != 0 {
			found = true
			examB = vB
			break
		}
	}

	//if no conflict is found make a single exam move (move examA to PeriodB)
	if found == false {
		solution[periodA] = append(solution[periodA][:exA], solution[periodA][exA+1:]...)
		i := sort.Search(len(solution[periodB]), func(i int) bool { return solution[periodB][i] >= examA })
		solution[periodB] = append(solution[periodB], 0)
		copy(solution[periodB][i+1:], solution[periodB][i:])
		solution[periodB][i] = examA
		return
	}

	//Get the Contents of the Periods
	lessonMap := make(map[int]map[int]bool)
	for _, per := range []int{periodA, periodB} {
		for _, v := range solution[per] {
			lessonMap[v] = make(map[int]bool)
			lessonMap[v][per] = false
		}
	}
	//find the chain
	lessonMap[examA][periodA] = true
	lessonMap[examB][periodB] = true
	next := []int{examA, examB}
	kempeChain(lessonMap, next, load)
	//Perform the exchange
	solution[periodA] = nil
	solution[periodB] = nil
	for k, v := range lessonMap {
		for key, value := range v {
			if (key == periodA && value == true) || (key == periodB && value == false) {
				solution[periodB] = append(solution[periodB], k)
			} else {
				solution[periodA] = append(solution[periodA], k)
			}
		}
	}

	sort.Ints(solution[periodA])
	sort.Ints(solution[periodB])
	return
}

func kempeChain(lessonMap map[int]map[int]bool, next []int, load *models.Loaded) {
	nextNext := []int{}
	for _, n := range next {
		var nPer int
		for nPeriod := range lessonMap[n] {
			nPer = nPeriod
		}
		for lesson, value := range lessonMap {
			if n != lesson {
				for period, v := range value {
					if v == false && period != nPer {
						if load.GraphNodes[lesson][n] > 0 {
							lessonMap[lesson][period] = true
							nextNext = append(nextNext, lesson)
						}
					}
				}
			}
		}
	}
	if len(nextNext) > 0 {
		kempeChain(lessonMap, nextNext, load)
	}
}
