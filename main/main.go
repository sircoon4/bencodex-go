package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/planetarium/bencodex-go"
)

func main() {
	var b []byte
	var err error

	//Avartar Example
	b, err = hex.DecodeString("6475373a747970655f69647531353a6372656174655f617661746172313175363a76616c7565736475333a65617269306575343a6861697269306575323a696431363aff371ccdebad6d439eaaa6228b9799f075353a696e64657869316575343a6c656e7369306575343a6e616d6575383a4a6f686e20446f6575343a7461696c6930656565")
	if err != nil {
		panic(err)
	}
	rv, err := bencodex.Decode(b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("value: %s\n type: %T\n", rv, rv)
		fmt.Println()
	}

	//SerializedPayload Example
	b, err = base64.StdEncoding.DecodeString("ZDE6UzcwOjBEAiBmJiuPcJWivuQlM9UTsWO9rl+YIegFAqOuD/ww27PcmwIgQ9XOjm4tUJS9N10PC6KFBgLXv0G8jhDFRbFMfi0Rd2gxOmFsZHU3OnR5cGVfaWR1MTY6aGFja19hbmRfc2xhc2gyMnU2OnZhbHVlc2R1MTI6YXBTdG9uZUNvdW50dTE6MHUxMzphdmF0YXJBZGRyZXNzMjA6nytBBv/0XCikHtGDYzFzi71nVM51ODpjb3N0dW1lc2xldTEwOmVxdWlwbWVudHNsMTY6RCeeIqAQWkq76wl93aAcMjE2Ok9UrTR1KllEqqkZnS1UT5YxNjoUTp5N7gmzQ6Mz38P3mc+iMTY6HM/tYYT6CEKGXdhq8jplPzE2OoJ+SHyYQbdLofs5WgQAJRYxNjqs166tESwlRJ/pT8XC1645MTY6svhM9XjLbUeHqr3JA/x1ZmV1NTpmb29kc2xldTI6aWQxNjpqQGwoVkUAR6gNm2dIlkCHdTE6cmxsdTE6MHU1OjMwMDAxZWV1NzpzdGFnZUlkdTM6MjEwdTE0OnRvdGFsUGxheUNvdW50dTE6MXU3OndvcmxkSWR1MTo1ZWVlMTpnMzI6RYIlDQ2jOwZ3moR10oPV3SEMaDubmZ100D+sT1j6a84xOmxpNGUxOm1sZHUxMzpkZWNpbWFsUGxhY2VzMToSdTc6bWludGVyc251Njp0aWNrZXJ1NDpNZWFkZWkxMDAwMDAwMDAwMDAwMDAwMDAwZWUxOm5pOTgzZTE6cDY1OgQJk5rs8nmCnnRsgbjJwEw0uMnel2SU/VB7enq6hCSfKHDLnen3eS1LZXqwW1d3FCw4MLnZ2vyFGii3hGOOmwM0MTpzMjA6LVWZmGl+p9oemMsRe3Mh1vmX8DExOnR1Mjc6MjAyNC0wNi0yNlQwNzozOToxNi44NzU1MDNaMTp1bGVl")
	if err != nil {
		panic(err)
	}
	rv, err = bencodex.Decode(b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("value: %s\n type: %T\n", rv, rv)
		fmt.Println()
	}
}