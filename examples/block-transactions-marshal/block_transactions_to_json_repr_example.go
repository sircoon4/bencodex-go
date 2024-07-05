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

	"github.com/joho/godotenv"
	"github.com/sircoon4/bencodex-go"
	"github.com/sircoon4/bencodex-go/util"
)

// Parse the serialized payload of a block transaction from the GraphQL query response
func blockTransactionsToJsonReprExample() {
	const dirPath = "bencodex_json_repr_datas"
	const filePath = "bencodex_json_repr_datas/bencodex_json_repr_data_%d.repr.json"
	const filePathForGlob = "bencodex_json_repr_datas/bencodex_json_repr_data_*.repr.json"
	const payloadDirPath = "bencodex_payloads"
	const payloadFilePath = "bencodex_payloads/bencodex_payload_%d.dat"
	const payloadFilePathForGlob = "bencodex_payloads/bencodex_payload_*.dat"

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	path9c := os.Getenv("9C_GRAPHQL_EXPLORER_API_URL")

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

	files, err := filepath.Glob(payloadFilePathForGlob)
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
	if err := os.MkdirAll(payloadDirPath, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	var serializedPayloadEncodedList [][]byte
	for i, transaction := range response.Data.BlockQuery.Blocks[0].Transactions {
		err = os.WriteFile(fmt.Sprintf(payloadFilePath, i), []byte(transaction.SerializedPayload), 0644)
		if err != nil {
			fmt.Println("Error writing serializedPayload data:", err)
			return
		}

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

	files, err = filepath.Glob(filePathForGlob)
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

		out, err := util.MarshalJsonRepr(serializedPayload)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		err = os.WriteFile(fmt.Sprintf(filePath, i), out, 0644)
		if err != nil {
			fmt.Println("Error writing Bencodex json data:", err)
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

		jsonData, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		data, err := util.UnmarshalJsonRepr(jsonData)
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
