package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sircoon4/bencodex-go/util"
)

func main() {
	const jsonFilePath = "../../../spec/testsuite/*.json"
	const jsonExcludeFilePath = "../../../spec/testsuite/*.repr.json"

	files, err := filepath.Glob(jsonFilePath)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	excluededFiles, err := filepath.Glob(jsonExcludeFilePath)
	if err != nil {
		fmt.Println("Error getting excluded files:", err)
	}
	// Exclude the files that are not including json map data
	files = difference(files, excluededFiles)
	for _, file := range files {
		fmt.Println("Test File:", file)

		jsonMapData, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading json file:", err)
			return
		}

		data, err := util.UnmarshalJsonMap(jsonMapData)
		if err != nil {
			fmt.Println("Error parsing json data:", err)
			return
		}
		fmt.Printf("Bencodex data: %v Type:%T\n", data, data)

		jsonMapData, err = util.MarshalJsonMap(data)
		if err != nil {
			fmt.Println("Error converting to json repr:", err)
			return
		}
		fmt.Println("Json Repr:", string(jsonMapData))

		data, err = util.UnmarshalJsonMap(jsonMapData)
		if err != nil {
			fmt.Println("Error parsing json data:", err)
			return
		}
		fmt.Printf("Bencodex data: %v Type:%T\n", data, data)

		fmt.Println()
	}
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// if string is not foundec, add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}
