package loadproblem

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//LoadProblem loads the problem into memory and creates the collision map and the cost map
func LoadProblem(load *models.Loaded) {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to read working directory!")
	}

	confPath := filepath.Join(currentDir, "datasets", "problems")
	files, err := ioutil.ReadDir(confPath)
	if err != nil {
		log.Fatal("Unable to read directory")
	}
	i := 1
	for _, file := range files {
		fmt.Println(strconv.Itoa(i) + " " + strings.TrimSuffix(file.Name(), ".stu"))
		i++
	}

	var selection int
	fmt.Scan(&selection)
	selection--
	if len(files) < selection || selection < 0 {
		fmt.Println("Wrong Selection")
	} else {
		OpenFile(load, selection, files, currentDir)
	}
}

//OpenFile opens a file
func OpenFile(load *models.Loaded, selection int, files []os.FileInfo, currentDir string) {

	load.Problem = ""
	load.Solution = ""
	load.SimoultaneousPairs = make(map[models.Pairs]int)
	filepath := filepath.Join(currentDir, "datasets", "problems", files[selection].Name())
	readFile(filepath, load)
	name := strings.TrimSpace(files[selection].Name())
	periods := 0
	findPeriods(name, &periods)
	for periods == 0 {
		fmt.Println("Periods not found, enter manually")
		var selection int
		fmt.Scan(&selection)
		periods = selection
	}
	load.Problem = files[selection].Name()
	load.Periods = periods
}

func readFile(path string, load *models.Loaded) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening problem")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal("Scanning error in problem scanning")
	}

	maxLesson := -1
	students := 0
	for scanner.Scan() {
		students++
		scanLine := strings.Split(scanner.Text(), " ")
		scanLineInt := []int{}
		for _, lesson := range scanLine {
			trim := strings.TrimSpace(lesson)
			if trim != "" {
				lessonInt, err := strconv.Atoi(trim)
				if err != nil {
					log.Fatal("Error reading problemfile")
				}
				scanLineInt = append(scanLineInt, lessonInt)

				if lessonInt > maxLesson {
					maxLesson = lessonInt
				}
			}
		}
		//sort.Ints(scanLineInt)
		for i := 0; i < len(scanLineInt); i++ {
			for j := i + 1; j < len(scanLineInt); j++ {
				load.SimoultaneousPairs[models.Pairs{X: scanLineInt[i], Y: scanLineInt[j]}]++
			}
		}
	}

	load.Lessons = maxLesson
	load.Students = students
	load.GraphNodes = make(map[int]map[int]int)

	for i := 1; i <= maxLesson; i++ {
		load.GraphNodes[i] = make(map[int]int)
	}

	conflicts := 0
	for i := 1; i <= maxLesson; i++ {
		for j := i + 1; j <= maxLesson; j++ {
			simultaneous := load.SimoultaneousPairs[models.Pairs{X: i, Y: j}]
			load.GraphNodes[i][j] = simultaneous
			load.GraphNodes[j][i] = simultaneous

			if simultaneous != 0 {
				conflicts++

			}
		}
	}
	density := float64(conflicts*2) / float64(maxLesson*maxLesson)
	fmt.Println("Density is:", math.Round(density*100)/100)
}

func findPeriods(s string, periods *int) {
	switch s {
	case "car-f-92.stu":
		*periods = 32
	case "car-s-91.stu":
		*periods = 35
	case "ear-f-83.stu":
		*periods = 24
	case "hec-s-92.stu":
		*periods = 18
	case "kfu-s-93.stu":
		*periods = 20
	case "lse-f-91.stu":
		*periods = 18
	case "pur-s-93.stu":
		*periods = 42
	case "rye-s-93.stu":
		*periods = 23
	case "sta-f-83.stu":
		*periods = 13
	case "tre-s-92.stu":
		*periods = 23
	case "uta-s-92.stu":
		*periods = 35
	case "ute-s-92.stu":
		*periods = 10
	case "yor-f-83.stu":
		*periods = 21
	}
}
