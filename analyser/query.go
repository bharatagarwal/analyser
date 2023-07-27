package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func query() {
	fmt.Println("Querying...")
	benchmarkStart := time.Now()

	// Checking for existing results file; removing if present
	if _, err := os.Stat("./retention_percentages.json"); err == nil {
		fmt.Printf("Results file exists. Deleting to query afresh...\n\n")
		err := os.Remove("./retention_percentages.json")
		if err != nil {
			fmt.Printf("error deleting database file: %s\n", err.Error())
			os.Exit(1)
		}
	}

	// Accessing database
	db, err := sql.Open("sqlite", "./records.db")

	if err != nil {
		fmt.Printf("error creating database file: %s\n", err.Error())
		os.Exit(1)
	}

	// Getting starting and ending months
	months := getMonthsInRange(db)

	// Getting unique users for first month
	var returningUsers []int64
	initialCount := firstMonthUniqueUsers(db, months[0])
	returningUsers = append(returningUsers, initialCount)

	// Getting unique returning users for subsequent months
	for _, month := range months[1:] {
		count := queryIntersection(db, months[0], month)
		returningUsers = append(returningUsers, count)
	}

	// Calculating percentage of returning users
	var percentages []int64

	for _, users := range returningUsers {
		percent := users * 100 / initialCount
		percentages = append(percentages, percent)
	}

	writeToJSON(months, percentages)

	fmt.Println("Querying complete")
	fmt.Printf("Duration: %v\n",
		time.Now().Sub(benchmarkStart))
}

func writeToJSON(months []string, percentages []int64) {
	type RetentionData struct {
		Months      []string `json:"months"`
		Percentages []int64  `json:"percentages"`
	}

	data := RetentionData{
		Months:      months,
		Percentages: percentages,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("error marshalling query data: %s\n", err.Error())
		os.Exit(1)
	}

	err = os.WriteFile("retention_percentages.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("error writing JSON file: %s\n", err.Error())
		os.Exit(1)
	}
}

func queryIntersection(db *sql.DB, initial string,
	current string) int64 {
	var count int64

	query := `
		SELECT COUNT(DISTINCT user_id)
		FROM inferences
		WHERE user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '` + initial + `'
		)
		AND user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '` + current + `'
		);`

	err := db.QueryRow(query).Scan(&count)

	if err != nil {
		fmt.Printf("error executing returning count query for month %s: %s\n",
			current, err.Error())
		os.Exit(1)
	}

	return count
}

func firstMonthUniqueUsers(db *sql.DB, month string) (count int64) {
	query := `SELECT count(DISTINCT user_id) 
		FROM inferences WHERE strftime('%Y-%m', timestamp) = '` + month + "';"

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Printf("error executing count query for month %s: %s\n",
			month, err.Error())
		os.Exit(1)
	}

	return count
}

func getMonthsInRange(db *sql.DB) []string {
	var months []string

	start, end := bookends(db)

	current := start

	for inRange(current, start, end) {
		monthYearString := fmt.Sprintf("%04d-%02d", current.Year(),
			current.Month())
		months = append(months, monthYearString)
		current = current.AddDate(0, 1, 0)
	}

	return months
}

func bookends(db *sql.DB) (time.Time, time.Time) {
	var start, end string
	_ = db.QueryRow("SELECT MIN(timestamp) FROM inferences").Scan(&start)
	_ = db.QueryRow("SELECT MAX(timestamp) FROM inferences").Scan(&end)

	startingTimestamp, _ := time.Parse("2006-01-02 15:04:05.999999-07:00",
		start)
	endingTimestamp, _ := time.Parse("2006-01-02 15:04:05.999999-07:00",
		end)

	return startingTimestamp, endingTimestamp
}

func inRange(current, start, end time.Time) bool {
	firstMonth := time.Date(start.Year(), start.Month(), 1,
		0,
		0, 0,
		0,
		time.UTC)
	monthAfterLast := time.Date(end.Year(), end.Month()+1, 1, 0,
		0, 0, 0, time.UTC) // ignoring 13th month edge case

	return current.After(firstMonth) && current.Before(monthAfterLast)
}