package bencodex

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBencodexEncode(t *testing.T) {
	testFiles, err := filepath.Glob("testsuite/*.json")
	if err != nil {
		t.Fatal(err)
	}
	excluededFiles, err := filepath.Glob("testsuite/*.repr.json")
	if err != nil {
		t.Fatal(err)
	}
	testFiles = difference(testFiles, excluededFiles)

	testResults, err := filepath.Glob("testsuite/*.dat")
	if err != nil {
		t.Fatal(err)
	}

	for i, file := range testFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			fmt.Println()
			fmt.Println("Test File:", file)

			// Read the test file
			jsonData, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}

			// Decode the YAML data
			var preData map[string]any
			err = json.Unmarshal(jsonData, &preData)
			if err != nil {
				t.Fatal(err)
			}
			data, err := parseJsonData(preData)
			if err != nil {
				t.Fatal(err)
			}

			// Encode the data
			encoded, err := Encode(data)
			if err != nil {
				t.Fatal(err)
			}

			// Read the test result file
			result, err := os.ReadFile(testResults[i])
			if err != nil {
				t.Fatal(err)
			}

			// Compare the encoded data with the test result
			assert.Equal(t, result, encoded)
		})
	}
}

func _TestBencodexDecode(t *testing.T) {
	testFiles, err := filepath.Glob("testsuite/*.dat")
	if err != nil {
		t.Fatal(err)
	}

	testResultFiles, err := filepath.Glob("testsuite/*.json")
	if err != nil {
		t.Fatal(err)
	}
	excluededResultFiles, err := filepath.Glob("testsuite/*.repr.json")
	if err != nil {
		t.Fatal(err)
	}
	testResults := difference(testResultFiles, excluededResultFiles)

	for i, file := range testFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			// Read the test file
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}

			// Decode the encoded data
			decoded, err := Decode(data)
			if err != nil {
				t.Fatal(err)
			}

			// Read the test file
			jsonData, err := os.ReadFile(testResults[i])
			if err != nil {
				t.Fatal(err)
			}

			// Decode the YAML data
			var preData map[string]any
			err = json.Unmarshal(jsonData, &preData)
			if err != nil {
				t.Fatal(err)
			}
			result, err := parseJsonData(preData)
			if err != nil {
				t.Fatal(err)
			}

			// Compare the original data with the decoded data
			assert.Equal(t, result, decoded)
		})
	}
}
