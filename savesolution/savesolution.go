package savesolution

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//SaveSolution saves a solution into the SavedSolutions folder
func SaveSolution(load *models.Loaded, solution models.Solution, score float64) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to read working directory!")
	}
	score = math.Round(score*100) / 100
	s := fmt.Sprintf("%.2f", score)
	fileName := strings.TrimSuffix(load.Problem, ".stu") + "(" + s + ").sol"
	filepath := filepath.Join(currentDir, "datasets", "savedSolutions", fileName)
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	for key, v := range solution.Solution {
		for _, value := range v {
			lesson := fmt.Sprintf("%04d", value) + "\t" + strconv.Itoa(key)

			buf := bytes.NewBufferString(lesson)
			fmt.Fprintln(f, buf)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("New Solution saved: ", filepath)
}
