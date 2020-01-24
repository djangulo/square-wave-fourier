package main

//go:generate go run main.go

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

const (
	imageHeight   = 260
	imageWidth    = 780
	chartDPI      = 300.0
	titleFontSize = 4.5
	axisFontSize  = 3.5
)

var (
	fontColor = drawing.Color{30, 30, 30, 255}
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
	var title string
	switch which {
	case "square":
		title = "Square wave"
		f = makeSquareWaveSeries
	default:
		if harmonicCount > 1 {
			title = fmt.Sprintf("Harmonics 1 - %d", harmonicCount)
		} else {
			title = "Harmonic 1"
		}
		f = makeHarmonicSeries
	}

	fontStyle := func(size float64) chart.Style {
		return chart.Style{
			FontSize:  size,
			FontColor: fontColor,
		}
	}

	xaxis := chart.XAxis{
		Name:           "Time (t)",
		NameStyle:      fontStyle(titleFontSize),
		ValueFormatter: xFormatter,
		Ticks:          makeXTicks(),
		TickPosition:   chart.TickPositionUnderTick,
		Style:          fontStyle(axisFontSize),
	}
	yaxis := chart.YAxis{
		Name:           "y=A*Sin(t*ω+phase)",
		NameStyle:      fontStyle(titleFontSize),
		ValueFormatter: yFormatter,
		Ticks:          makeYTicks(),
		Style:          fontStyle(axisFontSize),
	}

	graph := chart.Chart{
		Title:      title,
		TitleStyle: fontStyle(titleFontSize),
		Series:     f(),
		Width:      imageWidth,
		Height:     imageHeight,
		XAxis:      xaxis,
		YAxis:      yaxis,
		DPI:        chartDPI,
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
