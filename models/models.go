package models

import lru "github.com/hashicorp/golang-lru"

//Configuration struct stores all the settings imported at the start
type Configuration struct {
	Test              string `yaml:"test"`
	Shuffle           int    `yaml:"shuffle"`
	StepFactor        int    `yaml:"stepFactor"`
	RaiseStep         int    `yaml:"raiseStep"`
	ChangePeriods     int    `yaml:"changePeriods"`
	MassChangePeriods int    `yaml:"massChangePeriods"`
	MoveExam          int    `yaml:"moveExam"`
	ExchangeExams     int    `yaml:"exchangeExams"`
	MassExodus        int    `yaml:"massExodus"`
	MassMigration     int    `yaml:"massMigration"`
	CacheSize         int    `yaml:"cacheSize"`
	MaxSteps          int    `yaml:"maxSteps"`
	Restart           int    `yaml:"restart"`
}

//Loaded keeps the current problem, solution or both for easy reference
type Loaded struct {
	Problem            string
	Students           int
	Periods            int
	Lessons            int
	Solution           string
	SimoultaneousPairs map[Pairs]int
	GraphNodes         map[int]map[int]int
	Conf               Configuration
	MovesDist          MovesDistribution
	Cache              *lru.ARCCache
}

//MovesDistribution struct
type MovesDistribution struct {
	PossibilitySum    int
	ChangePeriods     int
	MoveExam          int
	ExchangeExams     int
	MassChangePeriods int
	MassExodus        int
	MassMigration     int
}

//Pairs struct
type Pairs struct {
	X int
	Y int
}

//Weight information for Greedy Coloring
type Weight struct {
	Lesson int
	Weight int
}

//Solution keeps a solution and its score
type Solution struct {
	Solution [][]int
	Score    int
}
