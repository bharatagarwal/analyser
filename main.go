package main

import (
	"fmt"
	"os"
)

func main() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		introduction()
		howToUse()
		return
	}

	command := arguments[0]

	if command == "ingest" {
		ingest()
	} else if command == "query" {
		query()
	} else if command == "chart" {
		chart()
	} else if command == "demo" {
		ingest()
		query()
		chart()
	} else {
		fmt.Println("Unknown command:", command, "\n")
		howToUse()
	}
}
