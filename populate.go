package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const TaintedTimestamp = "1970-01-01T00:00:00"

func convertTime(timestamp string) (time.Time, error) {
	format := "2006-01-02T15:04:05.000000"

	t, err := time.Parse(format, timestamp)
	if err != nil {
		return time.Time{}, err
	}

	return t.UTC(), nil
}

func convertBoolToBinary(value bool) int {
	if value == true {
		return 1
	}

	return 0
}

func populate(db *sql.DB) {
	inferences, _ := os.ReadFile("./inferences.json")
	var siteLog []Inference

	err := json.Unmarshal(inferences, &siteLog)
	if err != nil {
		fmt.Printf("error decoding json file: %s\n", err.Error())
		os.Exit(1)
	}

	schema := `
			CREATE TABLE inferences (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			is_anonymous BOOLEAN,
	        recipe TEXT,
	        run_id TEXT(64),
	        user_id TEXT(64),
	        timestamp DATETIME
		)`

	_, err = db.Exec(schema)
	// Ignoring result as information around
	// affected rows and last inserted id are not relevant

	if err != nil {
		fmt.Printf("error creating database schema: %s\n", err.Error())
		os.Exit(1)
	}

	insertStatement := `
			INSERT INTO inferences(
	            is_anonymous,
	            recipe,
	            run_id,
	            user_id,
	            timestamp
			)
			VALUES(?, ?, ?, ?, ?)
	`

	statement, err := db.Prepare(insertStatement)
	if err != nil {
		fmt.Printf("error preparing insertion statement: %s\n", err.Error())
		os.Exit(1)
	}

	for _, record := range siteLog {
		if record.Timestamp == TaintedTimestamp {
			continue
		}

		timestamp, err := convertTime(record.Timestamp)
		if err != nil {
			fmt.Printf("error parsing timezone: %s\n", err.Error())
			os.Exit(1)
		}

		anonymity := convertBoolToBinary(record.IsAnonymous)

		_, err = statement.Exec(
			anonymity,
			record.Recipe,
			record.RunID,
			record.UserID,
			timestamp,
		)

		if err != nil {
			fmt.Println(record)
			fmt.Printf("error inserting record: %s\n", err.Error())
			os.Exit(1)
		}
	}
}