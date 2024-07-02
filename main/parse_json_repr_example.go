package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sircoon4/bencodex-go/util"
)

func parseJsonReprExample() {
	const jsonFilePath = "../spec/testsuite/*.repr.json"

	files, err := filepath.Glob(jsonFilePath)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	for _, file := range files {
		fmt.Println("Test File:", file)

		jsonData, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading json file:", err)
			return
		}

		data, err := util.UnmarshalJsonRepr(jsonData)
		if err != nil {
			fmt.Println("Error parsing json data:", err)
			return
		}
		fmt.Printf("Bencodex data: %v Type:%T\n", data, data)

		jsonReprData, err := util.MarshalJsonRepr(data)
		if err != nil {
			fmt.Println("Error converting to json repr:", err)
			return
		}
		fmt.Println("Json Repr:", string(jsonReprData))

		data, err = util.UnmarshalJsonRepr(jsonReprData)
		if err != nil {
			fmt.Println("Error parsing json data:", err)
			return
		}
		fmt.Printf("Bencodex data: %v Type:%T\n", data, data)

		fmt.Println()
	}
}
