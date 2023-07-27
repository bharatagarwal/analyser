package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"os"
)

type Inference struct {
	IsAnonymous bool
	Recipe      string
	RunID       string `json:"run_id"`
	UserID      string `json:"user_id"`
	Timestamp   string
}

func main() {
	var populateFlag = flag.Bool("populate", false, "populate the database")
	var queryFlag = flag.Bool("query", false, "populate the database")
	var chartFlag = flag.Bool("chart", false, "populate the database")

	flag.Parse()

	db, err := sql.Open("sqlite", "./records.db")
	if err != nil {
		fmt.Printf("error opening database file: %s\n", err.Error())
		os.Exit(1)
	}

	if *populateFlag {
		populate(db)
	}

	if *queryFlag {
		query(db)
	}

	if *chartFlag {
		chart()
	}
}