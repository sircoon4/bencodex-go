package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sircoon4/bencodex-go"
	"github.com/sircoon4/bencodex-go/util"
)

func main() {
	const path9c = "https://9c-main-rpc-1.nine-chronicles.com/graphql/explorer"
	const jsonDirPath = "bencodex_json_datas"
	const jsonFilePath = "bencodex_json_datas/bencodex_json_data_%d.json"
	const jsonFilePathForGlob = "bencodex_json_datas/bencodex_json_data_%d.json"
	const yamlDirPath = "bencodex_yaml_datas"
	const yamlFilePath = "bencodex_yaml_datas/bencodex_yaml_data_%d.yaml"
	const yamlFilePathForGlob = "bencodex_yaml_datas/bencodex_yaml_data_*.yaml"

	// Make GraphQL query request
	query := `{
		blockQuery{
			blocks(desc: true, limit: 1) {
				transactions {
					serializedPayload
				}
			}
		}
	}`

	// Create the request body
	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	// Send the request
	resp, err := http.Post(path9c, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the response body
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	// Decode base64 encoded serialized payload
	var serializedPayloadEncodedList [][]byte
	for _, transaction := range response.Data.BlockQuery.Blocks[0].Transactions {
		serializedPayloadEncoded, err := base64.StdEncoding.DecodeString(transaction.SerializedPayload)
		if err != nil {
			fmt.Println("Error decoding serialized payload:", err)
			return
		}
		serializedPayloadEncodedList = append(serializedPayloadEncodedList, serializedPayloadEncoded)
	}

	// Dcode bencodex encoded serialized payload
	var serializedPayloadList []any
	for _, serializedPayloadEncoded := range serializedPayloadEncodedList {
		value, err := bencodex.Decode(serializedPayloadEncoded)
		if err != nil {
			fmt.Println("Error decoding bencodex value:", err)
			return
		}
		serializedPayloadList = append(serializedPayloadList, value)
	}

	// Delete existing files
	files, err := filepath.Glob(jsonFilePathForGlob)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	yamlFiles, err := filepath.Glob(yamlFilePathForGlob)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	files = append(files, yamlFiles...)
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
	}

	// Marshal bencodex value to JSON and write to file
	if err := os.Mkdir(jsonDirPath, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	for i, serializedPayload := range serializedPayloadList {
		out, err := util.MarshalJsonMap(serializedPayload)
		if err != nil {
			fmt.Println("Error marshalling json data:", err)
			return
		}

		err = os.WriteFile(fmt.Sprintf(jsonFilePath, i), out, 0644)
		if err != nil {
			fmt.Println("Error writing json data:", err)
			return
		}
	}

	// Marshal bencodex value to YAML and write to file
	if err := os.Mkdir(yamlDirPath, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	for i, serializedPayload := range serializedPayloadList {
		out, err := util.MarshalYaml(serializedPayload)
		if err != nil {
			fmt.Println("Error marshalling yaml data:", err)
			return
		}

		err = os.WriteFile(fmt.Sprintf(yamlFilePath, i), out, 0644)
		if err != nil {
			fmt.Println("Error writing yaml data:", err)
			return
		}
	}
}
