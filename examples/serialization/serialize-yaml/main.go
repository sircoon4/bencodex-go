package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sircoon4/bencodex-go/util"
)

func main() {
	const yamlFilePath = "../../../spec/testsuite/*.yaml"

	files, err := filepath.Glob(yamlFilePath)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	for _, file := range files {
		fmt.Println("Test File:", file)

		yamlData, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading json file:", err)
			return
		}

		data, err := util.UnmarshalYaml(yamlData)
		if err != nil {
			fmt.Println("Error parsing json data:", err)
			return
		}
		fmt.Printf("Bencodex data: %v Type:%T\n", data, data)

		yamlData, err = util.MarshalYaml(data)
		if err != nil {
			fmt.Println("Error converting to json repr:", err)
			return
		}
		fmt.Println("Json Repr:", string(yamlData))

		data, err = util.UnmarshalYaml(yamlData)
		if err != nil {
			fmt.Println("Error parsing json data:", err)
			return
		}
		fmt.Printf("Bencodex data: %v Type:%T\n", data, data)

		fmt.Println()
	}
}
