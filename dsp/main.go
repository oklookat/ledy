package dsp

import (
	"math"
	"math/cmplx"

	"github.com/oklookat/ledy/ears"
)

// Calculate the sum of amplitudes for all elements in the spectrum.
func AmplitudeMulti(spectrum []complex128) float64 {
	sum := 0.0
	for i := range spectrum {
		sum += AmplitudeSingle(spectrum[i])
	}
	return sum
}

// Return the amplitude of a single element in the FFT spectrum.
func AmplitudeSingle(fft complex128) float64 {
	return cmplx.Abs(fft)
}

func FrequencyRangeIndexes(fft []complex128, sampleRate ears.SampleRate, fromHz, toHz uint16) (int, int) {
	if len(fft) == 0 || sampleRate == 0 {
		return 0, 0
	}
	k := sampleRate / ears.SampleRate(len(fft))
	fromIndex := int(math.Round(float64(fromHz) / float64(k)))
	toIndex := int(math.Round(float64(toHz) / float64(k)))
	if fromIndex > toIndex {
		fromIndex, toIndex = toIndex, fromIndex
	}
	if fromIndex < 0 {
		fromIndex = 0
	}
	if toIndex >= len(fft) {
		toIndex = len(fft) - 1
	}
	return fromIndex, toIndex
}
