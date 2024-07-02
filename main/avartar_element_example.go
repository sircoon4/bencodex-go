package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sircoon4/bencodex-go"
)

// Avartar Element Example
func avartarElementExample() {
	b, err := hex.DecodeString("6475373a747970655f69647531353a6372656174655f617661746172313175363a76616c7565736475333a65617269306575343a6861697269306575323a696431363aff371ccdebad6d439eaaa6228b9799f075353a696e64657869316575343a6c656e7369306575343a6e616d6575383a4a6f686e20446f6575343a7461696c6930656565")
	if err != nil {
		panic(err)
	}
	rv, err := bencodex.Decode(b)
	if err != nil {
		fmt.Println("avartarElementExample: ", err)
	} else {
		fmt.Printf("value: %s\n type: %T\n", rv, rv)
		fmt.Println()
	}
}
