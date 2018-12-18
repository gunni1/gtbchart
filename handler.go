package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wcharczuk/go-chart"
	"log"
	"net/http"
	"time"
)

func main() {
	apiPort := flag.String("apiPort", "7000", "REST API Port")

	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/chart/series", TimeSeriesChartHandler)

	log.Fatal(http.ListenAndServe(":"+*apiPort, router))
}

func index(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "Hello")
}

func TimeSeriesChartHandler(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var chartDto TimeSeriesChartDto
	error := decoder.Decode(&chartDto)
	if error != nil {
		http.Error(response, error.Error(), http.StatusBadRequest)
	} else {
		graph := buildChart(chartDto)
		response.Header().Set("Content-Type", "image/png")
		graph.Render(chart.PNG, response)
	}
}

func buildChart(chartDto TimeSeriesChartDto) chart.Chart {

	series := []chart.Series{
		chart.TimeSeries{
			XValues: millisToTime(chartDto.XTimeMillis),
			YValues: chartDto.YValues,
		},
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      chartDto.XCaption,
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      chartDto.YCaption,
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: series,
	}
	return graph
}

func millisToTime(millis []int64) []time.Time {
	times := make([]time.Time, len(millis))
	for idx, milli := range millis {
		times[idx] = time.Unix(0, milli*int64(time.Millisecond))
	}
	return times
}

type TimeSeriesChartDto struct {
	XCaption    string    `json:"xCaption"`
	YCaption    string    `json:"yCaption"`
	XTimeMillis []int64   `json:"xTimeMillis"`
	YValues     []float64 `json:"yValues"`
}
