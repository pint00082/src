package checksolution

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//CheckSolution validates and evaluates a solution
func CheckSolution(load *models.Loaded) {
	if load.Problem == "" {
		fmt.Println("Please first load a Problem then check for a solution")
		return
	}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to read working directory!")
		return
	}

	confPath := filepath.Join(currentDir, "datasets", "solutions")
	files, err := ioutil.ReadDir(confPath)
	if err != nil {
		fmt.Println("Unable to read directory")
		return
	}
	i := 1
	for _, file := range files {
		fmt.Println(strconv.Itoa(i) + " " + strings.TrimSuffix(file.Name(), ".sol"))
		i++
	}

	var selection int
	fmt.Scan(&selection)
	selection--
	if len(files) < selection || selection < 0 {
		fmt.Println("Wrong Selection")
	} else {
		filepath := filepath.Join(currentDir, "datasets", "solutions", files[selection].Name())
		readFile(filepath, load)

	}
}

func readFile(path string, load *models.Loaded) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening solution")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Println("Scanning error in solution scanning")
		return
	}

	periods := make(map[int][]int)
	for scanner.Scan() {
		txt := strings.ReplaceAll(scanner.Text(), "\t", " ")
		scanLine := strings.Split(txt, " ")
		scanLineInt := []int{}
		for _, v := range scanLine {
			trim := strings.TrimSpace(v)
			vInt, err := strconv.Atoi(trim)
			if err != nil {
				fmt.Println("Error reading solutionFile")
				return
			}

			scanLineInt = append(scanLineInt, vInt)
		}
		if scanLineInt[1] < 0 || scanLineInt[1] > load.Periods-1 {
			fmt.Println("Error reading solutionFile")
			return

		}
		if scanLineInt[0] <= 0 || scanLineInt[0] > load.Lessons {
			fmt.Println("Error reading solutionFile")
			return
		}
		periods[scanLineInt[1]] = append(periods[scanLineInt[1]], scanLineInt[0])
	}
	//check that all lessons are present
	var valLessons []bool
	for i := 0; i <= load.Lessons; i++ {
		valLessons = append(valLessons, false)
	}
	for _, value := range periods {
		for _, v := range value {
			valLessons[v] = true
		}
	}
	validationOfLessons := true
	for i := 1; i <= load.Lessons; i++ {
		if valLessons[i] == false {
			validationOfLessons = false
		}
	}
	if validationOfLessons == false {
		fmt.Println("not all lessons are present")
		return
	}

	//validate constraints
	for _, value := range periods {
		for i := 0; i < len(value); i++ {
			for j := i; j < len(value); j++ {
				if load.GraphNodes[value[i]][value[j]] != 0 {
					fmt.Println("validation failed conflicts in schedule")
					return
				}
			}
		}
	}
	totalCost := CalculateCost(periods, load)
	score := float64(totalCost) / float64(load.Students)
	fmt.Println("Validated")
	fmt.Println("Score: ", math.Round(score*100)/100)
}

//CalculateCost finds the cost of a valid solution
func CalculateCost(periods map[int][]int, load *models.Loaded) int {
	totalCost := 0
	for i := 1; i < load.Periods; i++ {
		maxImpact := i - 5
		if maxImpact < 0 {
			maxImpact = 0
		}
		for j := i - 1; j >= maxImpact; j-- {
			periodCost := 0
			for _, vi := range periods[i] {
				for _, vj := range periods[j] {
					periodCost += load.GraphNodes[vi][vj]
				}
			}
			totalCost += periodCost * FindCostFactor(i, j)
		}

	}
	return totalCost
}

//FindCostFactor retruns the penalty factor
func FindCostFactor(x int, y int) int {
	diff := x - y
	if diff < 0 {
		diff *= -1
	}
	switch diff {
	case 1:
		return 16
	case 2:
		return 8
	case 3:
		return 4
	case 4:
		return 2
	case 5:
		return 1
	}
	return 0
}
