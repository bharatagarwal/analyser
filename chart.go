package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func chart() {
	fmt.Println("Charting...")
	benchmarkStart := time.Now()

	// Checking for existing chart file; removing if present
	if _, err := os.Stat("./chart.html"); err == nil {
		fmt.Printf("Chart page exists. Deleting to generate afresh...\n\n")
		err := os.Remove("./chart.html")
		if err != nil {
			fmt.Printf("error deleting chart page: %s\n", err.Error())
			os.Exit(1)
		}
	}

	// Opening retention stats
	file, err := os.Open("retention_percentages.json")
	if err != nil {
		fmt.Printf("error opening analysis json file: %s\n", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	// Reading from JSON file
	var data RetentionData

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("error decoding analysis json file: %s\n", err.Error())
		os.Exit(1)
	}

	// Generating chart page
	page := components.NewPage()
	page.AddCharts(monthlyRetention(data.Months, data.Percentages))

	f, err := os.Create("chart.html")
	if err != nil {
		panic(err)
	}

	// Writing to chart file
	err = page.Render(io.MultiWriter(f))
	if err != nil {
		return
	}

	fmt.Println("Chart generation complete. See chart.html")
	fmt.Printf("Duration: %v\n\n",
		time.Now().Sub(benchmarkStart))
}

func generateBarItems(values []int64) []opts.BarData {
	items := make([]opts.BarData, 0)

	for _, userCount := range values {
		items = append(items, opts.BarData{Value: userCount})
	}

	return items
}

func monthlyRetention(months []string, users []int64) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		// charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: false, Right: "80px"}),
		charts.WithTitleOpts(opts.Title{Title: "Monthly users retention chart" +
			""}),
	)

	bar.SetXAxis(months).
		AddSeries("Users", generateBarItems(users)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)

	return bar
}
