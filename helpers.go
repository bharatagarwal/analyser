package main

import (
	"fmt"
	"os"
	"time"
)

func introduction() {
	message := `Analyser is a tool for generating retention charts`

	fmt.Printf("%s\n\n", message)
}

func howToUse() {
	instructions := `Usage:

	./analyser <command>

The commands are:

	ingest      ingest "inferences.json" file in current directory
				into a SQLite database by the name of "records.db"

	query       query database for user retention related information
				and save into a JSON file named "retention_percentages.json"

	chart       chart monthly retention diagram into a file
				called "retention.png"
	
	demo		perform ingestion, querying and charting in sequence.


`

	fmt.Print(instructions)
}

func parseTime(timestamp string) time.Time {
	format := "2006-01-02T15:04:05.000000"

	t, err := time.Parse(format, timestamp)
	if err != nil {
		fmt.Printf("error parsing timestamp: %s\n", err)
		os.Exit(11)
	}

	return t.UTC()
}
