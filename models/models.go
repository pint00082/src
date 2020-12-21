package models

import lru "github.com/hashicorp/golang-lru"

//Configuration struct stores all the settings imported at the start
type Configuration struct {
	Test               string `yaml:"test"`
	Shuffle            int    `yaml:"shuffle"`
	MaxShuffle         int    `yaml:"maxShuffle"`
	RaiseShuffleFactor int    `yaml:"raiseShuffleFactor"`
	CacheSize          int    `yaml:"cacheSize"`
	Kempe              int    `yaml:"kempe"`
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
	Cache              *lru.ARCCache
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
