package main

import "time"

type Response struct {
	Data       Data       `json:"data"`
	Extensions Extensions `json:"extensions"`
}

type Data struct {
	BlockQuery BlockQuery `json:"blockQuery"`
}

type BlockQuery struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	SerializedPayload string `json:"serializedPayload"`
}

type Extensions struct {
	Tracing Tracing `json:"tracing"`
}

type Tracing struct {
	Version    int        `json:"version"`
	StartTime  time.Time  `json:"startTime"`
	EndTime    time.Time  `json:"endTime"`
	Duration   int        `json:"duration"`
	Parsing    Parsing    `json:"parsing"`
	Validation Validation `json:"validation"`
	Execution  Execution  `json:"execution"`
}

type Parsing struct {
	StartOffset int `json:"startOffset"`
	Duration    int `json:"duration"`
}

type Validation struct {
	StartOffSet int `json:"startOffset"`
	Duration    int `json:"duration"`
}

type Execution struct {
	Resolvers []any `json:"resolvers"`
}
