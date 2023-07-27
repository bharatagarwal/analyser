package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type Inference struct {
	IsAnonymous bool `json:"is_anonymous"`
	Recipe      string
	RunID       string `json:"run_id"`
	UserID      string `json:"user_id"`
	Timestamp   string
}

const FaultyTimestamp = "1970-01-01T00:00:00"

func ingest() {
	fmt.Println("Ingesting...")
	benchStart := time.Now()

	// Opening JSON file
	inferences, err := os.ReadFile("./inferences.json")

	if err != nil {
		fmt.Printf("error opening json file: %s\n", err.Error())
		os.Exit(1)
	}

	// Parsing inferences into memory
	var records []Inference

	err = json.Unmarshal(inferences, &records)

	if err != nil {
		fmt.Printf("error decoding json file: %s\n", err.Error())
		os.Exit(2)
	}

	// Creating database
	db, err := sql.Open("sqlite", "./records.db")

	if err != nil {
		fmt.Printf("error creating database file: %s\n", err.Error())
		os.Exit(3)
	}

	// Creating table
	schema := `CREATE TABLE inferences (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		is_anonymous BOOLEAN,
		recipe TEXT,
		run_id TEXT(64),
		user_id TEXT(64),
		timestamp DATETIME
	)`

	_, err = db.Exec(schema)

	if err != nil {
		fmt.Printf("error creating database schema: %s\n", err.Error())
		os.Exit(4)
	}

	// Preparing insertion template
	template := `INSERT INTO inferences(
		is_anonymous,
		recipe,
		run_id,
		user_id,
		timestamp
	) VALUES (?, ?, ?, ?, ?)`

	statement, err := db.Prepare(template)

	if err != nil {
		fmt.Printf("error preparing insertion statement: %s\n", err.Error())
		os.Exit(5)
	}

	// Inserting records
	for _, record := range records {
		if record.Timestamp == FaultyTimestamp {
			continue
		}

		_, err := statement.Exec(
			record.IsAnonymous,
			record.Recipe,
			record.RunID,
			record.UserID,
			parseTime(record.Timestamp),
		)

		if err != nil {
			fmt.Printf("error inserting record: %s\n", err.Error())
			os.Exit(6)
		}
	}

	benchEnd := time.Now()
	fmt.Printf("%Elapsed: %v\n", benchEnd.Sub(benchStart))
}
