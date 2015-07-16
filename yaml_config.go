package ukigumo

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func NewYAMLConfig() *Config {
	yamlFile := chooseEffectiveYAMLFile()
	if yamlFile != "" {
		buf, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		var c Config
		err = yaml.Unmarshal(buf, &c)
		return &c
	}

	log.Printf("There is no yaml file")
	return nil
}

func chooseEffectiveYAMLFile() string {
	files := []string{".ukigumo.yml", ".travis.yml"}
	for _, file := range files {
		_, err := os.Stat(file)
		if err == nil { // exists
			return file
		}
	}
	return ""
}
