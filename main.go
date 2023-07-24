package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

type inference struct {
	IsAnonymous bool
	Recipe      string
	RunID       string `json:"run_id"`
	UserID      string `json:"user_id"`
	Timestamp   string // Will convert to timestamp when inserting into SQL
}

func main() {
	inferences, _ := os.ReadFile("./inferences.json")
	var siteLog []inference

	err := json.Unmarshal(inferences, &siteLog)
	if err != nil {
		fmt.Printf("error decoding json file: %s\n", err.Error())
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", "./records.db")

	if err != nil {
		fmt.Printf("error creating database file: %s\n", err.Error())
		os.Exit(1)
	}

	schema := `CREATE TABLE inferences (
            is_anonymous BOOLEAN,
            recipe TEXT,
            run_id TEXT(64) PRIMARY KEY,
            user_id TEXT(64),
            timestamp DATETIME
         )`

	_, err = db.Exec(schema)

	if err != nil {
		fmt.Printf("error creating database schema: %s\n", err.Error())
		os.Exit(1)
	}
}