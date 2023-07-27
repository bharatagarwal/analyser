package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	benchmarkStart := time.Now()

	// Checking for existing database; removing if present
	if _, err := os.Stat("./records.db"); err == nil {
		fmt.Printf("Database file exists. Deleting to ingest afresh...\n\n")
		err := os.Remove("./records.db")
		if err != nil {
			fmt.Printf("error deleting database file: %s\n", err.Error())
			os.Exit(1)
		}
	}

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
		os.Exit(1)
	}

	// Creating database
	db, err := sql.Open("sqlite", "./records.db")

	if err != nil {
		fmt.Printf("error creating database file: %s\n", err.Error())
		os.Exit(1)
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
		os.Exit(1)
	}

	// Preparing insertion template
	const batchSize = 100
	const fieldsPerRow = 5
	var placeholders []string

	for i := 0; i < batchSize; i += 1 {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
	}

	templateStart := `INSERT INTO inferences(
		is_anonymous,
		recipe,
		run_id,
		user_id,
		timestamp
	) VALUES `

	template := templateStart + strings.Join(placeholders, ",")
	statement, err := db.Prepare(template)

	if err != nil {
		fmt.Printf("error preparing batch insertion statement: %s\n", err.Error())
		os.Exit(1)
	}

	// Inserting records in batches of 100
	var batchFields []any

	for i := 0; i < len(records); i += 1 {
		if records[i].Timestamp == FaultyTimestamp {
			continue
		}

		batchFields = append(batchFields,
			records[i].IsAnonymous,
			records[i].Recipe,
			records[i].RunID,
			records[i].UserID,
			parseTime(records[i].Timestamp),
		)

		if len(batchFields) == batchSize*fieldsPerRow {
			_, err := statement.Exec(batchFields...)
			if err != nil {
				fmt.Printf("error inserting record: %s\n", err.Error())
				os.Exit(1)
			}

			batchFields = []any{} // clearing the batch
		}
	}

	// Inserting remaining records
	if len(batchFields) > 0 {
		insertRemaining(batchFields, db, templateStart)
	}

	fmt.Println("Ingestion complete")
	fmt.Printf("Duration: %v\n",
		time.Now().Sub(benchmarkStart))
}

func insertRemaining(fields []any, db *sql.DB, templateStart string) {
	rawStatement := templateStart + "(?, ?, ?, ?, ?)"
	statement, err := db.Prepare(rawStatement)

	if err != nil {
		fmt.Printf("error preparing batch insertion statement: %s\n", err.Error())
		os.Exit(1)
	}

	for i := 0; i < len(fields); i += 5 {
		_, err := statement.Exec(fields...)
		if err != nil {
			fmt.Printf("error inserting record: %s\n", err.Error())
			os.Exit(1)
		}
	}
}