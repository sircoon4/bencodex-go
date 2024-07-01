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

	"github.com/planetarium/bencodex-go"
	"github.com/planetarium/bencodex-go/util"
)

func blockTransactionsToYamlExample() {
	const path9c = "https://9c-main-rpc-1.nine-chronicles.com/graphql/explorer"
	const dirPath = "bencodex_yaml_datas"
	const filePath = "bencodex_yaml_datas/bencodex_yaml_data_%d.yaml"
	const filePathForGlob = "bencodex_yaml_datas/bencodex_yaml_data_*.yaml"

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

	// Print the response body
	fmt.Printf("%#v\n", response)
	fmt.Println()

	var serializedPayloadEncodedList [][]byte
	for _, transaction := range response.Data.BlockQuery.Blocks[0].Transactions {
		serializedPayloadEncoded, err := base64.StdEncoding.DecodeString(transaction.SerializedPayload)
		if err != nil {
			fmt.Println("Error decoding serialized payload:", err)
			return
		}
		serializedPayloadEncodedList = append(serializedPayloadEncodedList, serializedPayloadEncoded)
	}

	var serializedPayloadList []any
	for _, serializedPayloadEncoded := range serializedPayloadEncodedList {
		value, err := bencodex.Decode(serializedPayloadEncoded)
		if err != nil {
			fmt.Println("Error decoding bencodex value:", err)
			return
		}
		serializedPayloadList = append(serializedPayloadList, value)
	}

	files, err := filepath.Glob(filePathForGlob)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
	}

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	for i, serializedPayload := range serializedPayloadList {
		out, err := util.MarshalYaml(serializedPayload)
		if err != nil {
			fmt.Println("Error marshalling Bencodex yaml data:", err)
			return
		}

		err = os.WriteFile(fmt.Sprintf(filePath, i), out, 0644)
		if err != nil {
			fmt.Println("Error writing Bencodex yaml data:", err)
			return
		}
	}

	files, err = filepath.Glob(filePathForGlob)
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	for _, file := range files {
		i := 0
		_, err := fmt.Sscanf(file, filePath, &i)
		if err != nil {
			fmt.Println("Error extracting number from file name:", err)
			return
		}

		yamlData, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		data, err := util.UnmarshalYaml(yamlData)
		if err != nil {
			fmt.Printf("Error parsing Bencodex map data: %v", err)
		}

		// Encode the data
		encoded, err := bencodex.Encode(data)
		if err != nil {
			fmt.Println("Error encoding data:", err)
		}

		serializedPayloadEncoded := serializedPayloadEncodedList[i]
		if !bytes.Equal(encoded, serializedPayloadEncoded) {
			fmt.Println("Encoded data does not match serialized payload")
		}
	}
}
