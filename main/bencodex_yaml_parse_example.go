package main

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const YamlDataFilesPath = "../spec/testsuite/*.yaml"
const YamlDataFilesPathExtended = "yaml_datas/*.yaml"

func bencodexYamlParseExample() {
	yamlFiles, err := filepath.Glob(YamlDataFilesPath)
	if err != nil {
		panic(err)
	}
	yamlFilesExtended, err := filepath.Glob(YamlDataFilesPathExtended)
	if err != nil {
		panic(err)
	}
	yamlFiles = append(yamlFiles, yamlFilesExtended...)
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

		out, err := yaml.Marshal(obj)
		if err != nil {
			panic(err)
		}
		fmt.Println("Yaml Data:", string(out))
	}

	testMap := map[string]any{"!!binary \"aGVsbG8=\"": "hello"}
	out, err := yaml.Marshal(testMap)
	if err != nil {
		panic(err)
	}
	fmt.Println("Yaml Data:", string(out))

	testBigInt := new(big.Int)
	testBigInt.SetString("1234567890123456789012345678901234567890", 10)
	out, err = yaml.Marshal(testBigInt)
	if err != nil {
		panic(err)
	}
	fmt.Println("Yaml Data:", string(out))

	testBigIntString := "\"BencodexBigIntInspector\"1234567890123456789012345678901234567890"
	out, err = yaml.Marshal(testBigIntString)
	if err != nil {
		panic(err)
	}
	fmt.Println("Yaml Data:", string(out))
}
