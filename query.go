package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

var months []string
var uniqueMonthlyUsers []int64
var returningUsers []int64
var percentages []int64

func dateInMonthRange(d time.Time, startingGoTime time.Time,
	endingGoTime time.Time,
) bool {
	current := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
	start := time.Date(startingGoTime.Year(), startingGoTime.Month(), 1,
		0,
		0, 0,
		0,
		time.UTC)
	end := time.Date(endingGoTime.Year(), endingGoTime.Month(), 1, 0,
		0, 0, 0, time.UTC)

	return (current.After(start) || current.Equal(start)) &&
		(current.Before(end) || current.Equal(end))
}

func query(db *sql.DB) {
	// programmatically generate months from dec 2022 to june 2023,
	// using the date/time libraries in the standard library
	// monthlyCount := make(map[time.Month]int64)

	d := startingGoTime

	for dateInMonthRange(d, startingGoTime, endingGoTime) {
		monthYearString := fmt.Sprintf("%04d-%02d", d.Year(),
			d.Month())
		months = append(months, monthYearString)
		d = d.AddDate(0, 1, 0)
	}

	for _, month := range months {
		var count int64

		// Prepare the SQL with the month
		sqlQuery := fmt.Sprintf("SELECT count("+
			"DISTINCT user_id) FROM inferences WHERE strftime('%%Y-%%m', timestamp) = '%s';", month)

		err := db.QueryRow(sqlQuery).Scan(&count)

		if err != nil {
			fmt.Printf("error executing count query for month %s: %s\n",
				month, err.Error())
			os.Exit(1)
		}

		uniqueMonthlyUsers = append(uniqueMonthlyUsers, count)
	}

	returningUsers = []int64{79, 24, 12, 13, 11, 13, 7}

	/*
		79

		SELECT COUNT(DISTINCT user_id), recipe
		FROM inferences
		WHERE user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2022-12'
		)
		AND user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2023-01'
		)
		GROUP BY recipe;

		24

		SELECT COUNT(DISTINCT user_id) AS unique_users
		FROM inferences
		WHERE user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2022-12'
		)
		AND user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2023-02'
		);

		12

		SELECT COUNT(DISTINCT user_id) AS unique_users
		FROM inferences
		WHERE user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2022-12'
		)
		AND user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2023-03'
		);
		13


		SELECT COUNT(DISTINCT user_id) AS unique_users
		FROM inferences
		WHERE user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2022-12'
		)
		AND user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2023-04'
		);
		11


		SELECT COUNT(DISTINCT user_id) AS unique_users
		FROM inferences
		WHERE user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2022-12'
		)
		AND user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2023-05'
		);
		13

		SELECT COUNT(DISTINCT user_id) AS unique_users
		FROM inferences
		WHERE user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2022-12'
		)
		AND user_id IN (
		  SELECT user_id
		  FROM inferences
		  WHERE strftime('%Y-%m', timestamp) = '2023-06'
		);
		7
	*/

	var percentages []int64

	initial := returningUsers[0]

	for _, users := range returningUsers {
		percent := users * 100 / initial
		percentages = append(percentages, percent)
	}

	fmt.Println(months)
	fmt.Println(returningUsers)
}