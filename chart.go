package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
)

func chart() {
	initial := returningUsers[0]

	for _, users := range returningUsers {
		percent := users * 100 / initial
		percentages = append(percentages, percent)
	}

	page := components.NewPage()
	page.AddCharts(monthlyRetention(months, percentages))

	f, err := os.Create("chart.html")
	if err != nil {
		panic(err)
	}

	err = page.Render(io.MultiWriter(f))
	if err != nil {
		return
	}
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