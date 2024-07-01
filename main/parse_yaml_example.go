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
		fmt.Println("Test File:", file)

		yamlData, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading yaml file:", err)
			return
		}

		bencodexData, err := util.UnmarshalYaml(yamlData)
		if err != nil {
			fmt.Println("Error parsing yaml data:", err)
			return
		}

		fmt.Println("Bencodex Data:", bencodexData)

		// Convert bencodex data to yaml
		yamlData, err = util.MarshalYaml(bencodexData)
		if err != nil {
			fmt.Println("Error converting bencodex data to yaml:", err)
			return
		}
		fmt.Println("Yaml Data:", string(yamlData))

		bencodexData, err = util.UnmarshalYaml(yamlData)
		if err != nil {
			fmt.Println("Error parsing yaml data:", err)
			return
		}

		fmt.Println("Bencodex Data:", bencodexData)

		fmt.Println()
	}
}
