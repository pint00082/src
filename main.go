package main

import (
	"checksolution"
	"configuration"
	"fmt"
	"loadproblem"
	"log"
	"models"
	"solveall"
	"solveproblem"

	lru "github.com/hashicorp/golang-lru"
)

func main() {
	load := newLoad()
	configuration.ReadConfiguration(&load.Conf)
	//Grab Memory for Cache
	cache, err := lru.NewARC(load.Conf.CacheSize)
	if err != nil {
		log.Fatal("Unable to create cache!")
	}
	load.Cache = cache

	//Create the UI
	for {
		fmt.Println()
		fmt.Println("Selection:")
		fmt.Println("Selected Problem:", load.Problem)
		fmt.Println("Students:", load.Students)
		fmt.Println("Lessons:", load.Lessons)
		fmt.Println("Periods:", load.Periods)
		fmt.Println()
		fmt.Println("Choose Function:")
		fmt.Println("1. Load Problem")
		fmt.Println("2. Check Solution")
		fmt.Println("3. Solve Problem")
		fmt.Println("4. Mass Solve All Problems")
		fmt.Println("0. EXIT")

		var selection int
		fmt.Scan(&selection)

		switch selection {
		case 1:
			loadproblem.LoadProblem(load)
		case 2:
			checksolution.CheckSolution(load)
		case 3:
			var seconds int
			fmt.Println("Seconds to try?")
			fmt.Scan(&seconds)
			if seconds > 0 {
				if load.Problem != "" {
					solveproblem.SolveProblem(load, seconds)
				} else {
					fmt.Println("First select a Problem")
				}
			} else {
				fmt.Println("Time must be greater than 0")
			}

		case 4:
			var seconds int
			fmt.Println("Seconds to try?")
			fmt.Scan(&seconds)
			if seconds > 0 {
				solveall.SolveAll(load, seconds)
			} else {
				fmt.Println("Time must be greater than 0")
			}
		case 0:
			return
		default:
			fmt.Println("Wrong Choice!")
		}
	}
}

func newLoad() *models.Loaded {
	var l models.Loaded
	l.Problem = ""
	l.Solution = ""
	l.SimoultaneousPairs = make(map[models.Pairs]int)

	return &l
}
