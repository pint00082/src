package move

import (
	"math/rand"
	"models"
)

//ChangePeriods exchanges two periods
func ChangePeriods(load *models.Loaded, climber [][]int, random *rand.Rand) {

	periodA := 0
	periodB := 0
	//choosePeriods to exchange
	for periodA == periodB {
		periodA = random.Intn(load.Periods)
		periodB = random.Intn(load.Periods)
	}
	//exchangeSequence
	if periodA != periodB {
		climber[periodA], climber[periodB] = climber[periodB], climber[periodA]
	}
}

//MassChangePeriods exchanges all periods randomly
func MassChangePeriods(load *models.Loaded, climber [][]int, random *rand.Rand) {
	shuffledPeriods := random.Perm(load.Periods)
	for i := 0; i < load.Periods; i += 2 {
		climber[shuffledPeriods[i]], climber[shuffledPeriods[i+1]] = climber[shuffledPeriods[i+1]], climber[shuffledPeriods[i]]
	}
}

//ExamMove will try to move a single exam into a suitable period
func ExamMove(load *models.Loaded, climber [][]int, random *rand.Rand) {

	period := random.Intn(load.Periods)
	if len(climber[period]) == 0 {
		return
	}
	exam := random.Intn(len(climber[period]))
	examNumber := climber[period][exam]
	moveExamToRandomPeriod(examNumber, period, load, climber, random)
}

//MassMigration tries to move one random exam from each period
func MassMigration(load *models.Loaded, climber [][]int, random *rand.Rand) {

	periods := random.Perm(load.Periods)
	for _, period := range periods {
		if len(climber[period]) == 0 {
			continue
		}
		exam := random.Intn(len(climber[period]))
		examNumber := climber[period][exam]
		moveExamToRandomPeriod(examNumber, period, load, climber, random)
	}
}

//MassExodus tries to eject from a random period all of it's lessons in random order
func MassExodus(load *models.Loaded, climber [][]int, random *rand.Rand) {
	period := random.Intn(load.Periods)
	for _, examNumber := range climber[period] {
		moveExamToRandomPeriod(examNumber, period, load, climber, random)
	}
}

//ExchangeExams tries to change two exams between two periods
func ExchangeExams(load *models.Loaded, climber [][]int, random *rand.Rand) {
	periodA := 0
	periodB := 0
	//choosePeriods to exchange
	for periodA == periodB {
		periodA = random.Intn(load.Periods)
		periodB = random.Intn(load.Periods)
	}

	//check exams from periodA that can go to periodB
	canMoveA := []int{}
	for _, examA := range climber[periodA] {
		cnt := 0
		for _, examB := range climber[periodB] {
			if load.GraphNodes[examA][examB] != 0 {
				cnt++
			}
		}
		if cnt == 1 {
			canMoveA = append(canMoveA, examA)
		}
	}
	//check exams from periodB that can go to periodA
	canMoveB := []int{}
	for _, examB := range climber[periodB] {
		cnt := 0
		for _, examA := range climber[periodA] {
			if load.GraphNodes[examA][examB] != 0 {
				cnt++
			}
		}
		if cnt == 1 {
			canMoveB = append(canMoveB, examB)
		}
	}

	//Check for an eligible move
	random.Shuffle(len(canMoveA), func(i, j int) { canMoveA[i], canMoveA[j] = canMoveA[j], canMoveA[i] })
	random.Shuffle(len(canMoveB), func(i, j int) { canMoveB[i], canMoveB[j] = canMoveB[j], canMoveB[i] })

	for _, examA := range canMoveA {
		for _, examB := range canMoveB {
			if load.GraphNodes[examA][examB] != 0 {
				removeFromPeriod(climber[periodA], examA)
				climber[periodB] = append(climber[periodB], examA)
				removeFromPeriod(climber[periodB], examB)
				climber[periodA] = append(climber[periodA], examB)
				return
			}
		}
	}

}

func moveExamToRandomPeriod(examNumber int, period int, load *models.Loaded, climber [][]int, random *rand.Rand) {
	permutation := random.Perm(load.Periods)
	for _, v := range permutation {
		if v != period {
			canJoin := true
			for _, value := range climber[v] {
				if load.GraphNodes[examNumber][value] != 0 {
					canJoin = false
					break
				}
			}
			if len(climber[v]) == 0 {
				canJoin = true
			}
			if canJoin {
				removeFromPeriod(climber[period], examNumber)
				climber[v] = append(climber[v], examNumber)
				return
			}
		}
	}
}
func removeFromPeriod(period []int, exam int) {
	var i int
	for k, v := range period {
		if exam == v {
			i = k
		}
	}
	period = append(period[:i], period[i+1:]...)
}
