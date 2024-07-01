package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/planetarium/bencodex-go/util"
)

func parseYamlExample() {
	const yamlFilePath = "../spec/testsuite/*.yaml"

	files, err := filepath.Glob(yamlFilePath)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}

	for _, file := range files {
		fmt.Println()
		fmt.Println("Test File:", file)

		yamlData, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading yaml file:", err)
			return
		}

		bencodexData, err := util.ParseYamlData(yamlData)
		if err != nil {
			fmt.Println("Error parsing yaml data:", err)
			return
		}

		fmt.Println("Bencodex Data:", bencodexData)
	}
}
