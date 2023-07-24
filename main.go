package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ISO8601Time time.Time

func (*ISO8601Time) UnmarshalJSON() {

}

type record struct {
	Anonymity bool        `json:"is_anonymous"`
	Recipe    string      `json:"recipe"`
	RunID     string      `json:"run_id"`
	UserID    string      `json:"user_id"`
	GivenTime ISO8601Time `json:"timestamp"`
}

func main() {
	minimal, _ := os.ReadFile("minimal.json") // returns a byte array

	var siteLog []record

	err := json.Unmarshal(minimal, &siteLog)
	if err != nil {
		fmt.Printf("error decoding json file: %s\n", err.Error())
		os.Exit(1)
	}

	for _, log := range siteLog {
		fmt.Printf("%#v\n", log)
	}

	time.Parse()
}