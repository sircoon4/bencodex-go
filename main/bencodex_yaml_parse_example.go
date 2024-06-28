package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const YamlDataFilesPath = "../spec/testsuite/*.yaml"

func bencodexYamlParseExample() {
	yamlFiles, err := filepath.Glob(YamlDataFilesPath)
	if err != nil {
		panic(err)
	}
	for _, file := range yamlFiles {
		fmt.Println()
		fmt.Println("Yaml File:", file)

		var obj any

		yamlData, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlData, &obj)
		if err != nil {
			panic(err)
		}
		fmt.Printf("value: %v\n type: %T\n", obj, obj)
	}
}
