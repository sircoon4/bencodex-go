package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/planetarium/bencodex-go"
)

// Parse the serialized payload of a block transaction from the GraphQL query response
// Get response from https://9c-main-rpc-1.nine-chronicles.com/graphql/explorer
func blockTransactionsParseExample() {
	const path9c = "https://9c-main-rpc-1.nine-chronicles.com/graphql/explorer"

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

	var serializedPayloadList []any
	for _, transaction := range response.Data.BlockQuery.Blocks[0].Transactions {
		d, err := base64.StdEncoding.DecodeString(transaction.SerializedPayload)
		if err != nil {
			fmt.Println("Error decoding serialized payload:", err)
			return
		}
		value, err := bencodex.Decode(d)
		if err != nil {
			fmt.Println("Error decoding bencodex value:", err)
			return
		}
		serializedPayloadList = append(serializedPayloadList, value)
	}
	fmt.Println("Serialized Payload List:", serializedPayloadList)
}
