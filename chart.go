package main

//go:generate go run main.go

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"

	chart "github.com/wcharczuk/go-chart"
)

const (
	imageHeight = 260
	imageWidth  = 780
)

func yFormatter(v interface{}) string {
	x := v.(float64)
	return fmt.Sprintf("%.3f", x)
}

func xFormatter(v interface{}) string {
	x := v.(float64)
	return fmt.Sprintf("%.3f", x)
}

func makeHarmonicSeries() []chart.Series {
	series := make([]chart.Series, 0)
	for _, h := range harmonics {
		s := chart.ContinuousSeries{
			Name:            h.Details(),
			XValues:         timeArray,
			YValues:         h.Values,
			XValueFormatter: xFormatter,
			YValueFormatter: yFormatter,
		}
		series = append(series, s)
	}
	return series
}

func makeSquareWaveSeries() []chart.Series {
	series := make([]chart.Series, 0)
	s := chart.ContinuousSeries{
		Name:            "Sum of harmonics",
		XValues:         timeArray,
		YValues:         harmonicSum,
		XValueFormatter: xFormatter,
		YValueFormatter: yFormatter,
	}
	series = append(series, s)
	return series
}

func makeXTicks() []chart.Tick {
	ticks := make([]chart.Tick, 0)
	for i := 0; i < 11; i++ {
		t := chart.Tick{
			Value: float64(i) / 10.0,
			Label: xFormatter(float64(i) / 10.0),
		}
		ticks = append(ticks, t)
	}
	return ticks
}

func makeYTicks() []chart.Tick {
	ticks := make([]chart.Tick, 0)
	for i := -1.5; i < 1.6; i += 0.5 {
		ticks = append(ticks, chart.Tick{Value: i, Label: yFormatter(i)})
	}
	return ticks
}

func getChart(which string) string {
	var f func() []chart.Series
	switch which {
	case "square":
		f = makeSquareWaveSeries
	default:
		f = makeHarmonicSeries
	}

	xaxis := chart.XAxis{
		Name:           "Time (t)",
		ValueFormatter: xFormatter,
		Ticks:          makeXTicks(),
		TickPosition:   chart.TickPositionUnderTick,
	}
	yaxis := chart.YAxis{
		Name:           "y=A*Sin(t*ω+phase)",
		ValueFormatter: yFormatter,
		Ticks:          makeYTicks(),
	}

	graph := chart.Chart{
		Series: f(),
		Width:  imageWidth,
		Height: imageHeight,
		XAxis:  xaxis,
		YAxis:  yaxis,
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Fatal(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return base64Str

}

func getSquareWave() string {
	return getChart("square")
}

func getHarmonics() string {
	return getChart("harmonics")
}
