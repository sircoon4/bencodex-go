package bencodex

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/planetarium/bencodex-go/util"
	"github.com/stretchr/testify/assert"
)

const encodedDataFilesPath = "spec/testsuite/*.dat"
const decodedDataFilesPath = "spec/testsuite/*.json"
const excludedDataFilesPath = "spec/testsuite/*.repr.json"

func TestBencodexEncode(t *testing.T) {
	testFiles, err := filepath.Glob(decodedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}
	excluededFiles, err := filepath.Glob(excludedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}
	// Exclude the files that are not for encoding
	testFiles = difference(testFiles, excluededFiles)

	testResults, err := filepath.Glob(encodedDataFilesPath)
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

			// Decode the Json data
			var preData map[string]any
			err = json.Unmarshal(jsonData, &preData)
			if err != nil {
				t.Fatal(err)
			}
			data, err := util.ParseBencodexJasonMapData(preData)
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

func TestBencodexDecode(t *testing.T) {
	testFiles, err := filepath.Glob(encodedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}

	testResultFiles, err := filepath.Glob(decodedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}
	excluededResultFiles, err := filepath.Glob(excludedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}
	// Exclude the files that are not for decoding results
	testResults := difference(testResultFiles, excluededResultFiles)

	for i, file := range testFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			fmt.Println()
			fmt.Println("Test File:", file)

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

			// Decode the Json data
			var preData map[string]any
			err = json.Unmarshal(jsonData, &preData)
			if err != nil {
				t.Fatal(err)
			}
			result, err := util.ParseBencodexJasonMapData(preData)
			if err != nil {
				t.Fatal(err)
			}

			// Compare the original data with the decoded data
			customizedAssertEqual(t, result, decoded)
		})
	}
}
