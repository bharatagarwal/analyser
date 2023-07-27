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
		fmt.Printf("\n------\n")
		fmt.Println("Check for stages, and give instructions.")
		fmt.Println("If all stages are covered, direct towards output and notes.")
		os.Exit(0)
	}

	command := arguments[0]

	if command == "ingest" {
		ingest()
	} else if command == "query" {
		query()
	} else if command == "chart" {
		chart()
	} else {
		fmt.Println("Unknown command:", command, "\n")
		howToUse()
	}
}
