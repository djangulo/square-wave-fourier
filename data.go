package main

import (
	"fmt"
	"math"
)

const (
	valueLength = 1000
)

var (
	harmonicCount int = 10
	harmonics     []*Harmonic
	baseFreq      int       = 5
	tInterval     float64   = 0.001
	timeArray     []float64 = make([]float64, valueLength)
	harmonicSum             = make([]float64, valueLength)
)

func init() {
	buildTimeArray()
	buildHarmonics()
	sumHarmonics()
}

func buildTimeArray() {
	timeArray[0] = 0.0
	for i := 1; i < valueLength; i++ {
		timeArray[i] = timeArray[i-1] + tInterval
	}
}

func buildHarmonics() {
	harmonics = make([]*Harmonic, harmonicCount)
	for i := 1; i <= harmonicCount; i++ {
		harmonics[i-1] = NewHarmonic(i)
	}
}

func sumHarmonics() {
	harmonicSum = make([]float64, valueLength)
	for i := 0; i < valueLength; i++ {
		for _, h := range harmonics {
			harmonicSum[i] += h.Values[i]
		}
	}
}

type Harmonic struct {
	No        int
	Amplitude float64
	Frequency float64
	Phase     float64
	Values    []float64
}

func (h *Harmonic) String() string {
	return fmt.Sprintf("Harmonic %d", h.No)
}

func (h *Harmonic) Details() string {
	v := len(harmonics)
	if v < 10 {
		return fmt.Sprintf("Harmonic %d: A=%.4f; f=%.1f; ω=%2.4f; Phase: %.1f",
			h.No, h.Amplitude, h.Frequency, h.AngularFrequency(), h.Phase)
	} else if v < 100 {
		return fmt.Sprintf("Harmonic %2d: A=%.4f; f=%.1f; ω=%2.4f; Phase: %.1f",
			h.No, h.Amplitude, h.Frequency, h.AngularFrequency(), h.Phase)
	} else {
		return fmt.Sprintf("Harmonic %3d: A=%.4f; f=%.1f; ω=%2.4f; Phase: %.1f",
			h.No, h.Amplitude, h.Frequency, h.AngularFrequency(), h.Phase)
	}
}

func NewHarmonic(number int) *Harmonic {
	h := &Harmonic{
		No:        number,
		Amplitude: getAmplitude(number),
		Frequency: getFrequency(number),
		Phase:     getPhase(number),
		Values:    make([]float64, len(timeArray)),
	}

	h.populateValues()
	return h
}

func (h *Harmonic) populateValues() {
	for i, t := range timeArray {
		h.Values[i] = h.ValueAt(t)
	}
}

func (h *Harmonic) ValueAt(t float64) float64 {
	return h.Amplitude * math.Sin((t*h.AngularFrequency())+h.Phase)
}

func getPhase(i int) float64 {
	return 0
}

func getFrequency(i int) float64 {
	return float64(i * baseFreq)
}

func (h *Harmonic) AngularFrequency() float64 {
	return 2.0 * math.Pi * h.Frequency
}

func getAmplitude(i int) float64 {
	if i%2 != 0 {
		return 4.0 / (float64(i) * math.Pi)
	}
	return 0.0
}
