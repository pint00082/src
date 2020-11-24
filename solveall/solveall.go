package solveall

import (
	"io/ioutil"
	"loadproblem"
	"log"
	"models"
	"os"
	"path/filepath"
	"solveproblem"
)

//SolveAll tries to solve all problems id dataset/problems
func SolveAll(load *models.Loaded, seconds int) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to read working directory!")
	}

	confPath := filepath.Join(currentDir, "datasets", "problems")
	files, err := ioutil.ReadDir(confPath)
	if err != nil {
		log.Fatal("Unable to read directory")
	}

	for l := 0; l < len(files); l++ {
		loadproblem.OpenFile(load, l, files, currentDir)
		solveproblem.SolveProblem(load, seconds)
	}
}
