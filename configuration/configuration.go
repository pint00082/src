package configuration

import (
	"io/ioutil"
	"log"
	"models"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//ReadConfiguration read the Configuration file  conf.yaml in YAML format located in the same folder as the application
func ReadConfiguration(conf *models.Configuration) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to read working directory!")
	}

	confPath := filepath.Join(currentDir, "conf.yaml")
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal("Unable to read configuration file!")
	}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("Unable to unmarshal conf.yaml")
	}
}
