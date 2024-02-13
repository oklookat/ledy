package effect

import (
	"math"

	"github.com/oklookat/ledy/dsp"
	"github.com/oklookat/ledy/ears"
)

const LedsCount = 240

func NewLEDS() LEDS {
	var res LEDS
	for i := 0; i < len(res); i++ {
		res[i] = &RGB{}
	}
	return res
}

type LEDS [LedsCount]*RGB

type RGB struct {
	R, G, B uint8
}

const _maxAmplitude = 50.0

func changeBrightness(color *RGB, percents float64) {
	if percents > 100.0 {
		percents = 100.0
	} else if percents < 0.0 {
		percents = 0.0
	}
	factor := percents / 100.0
	if factor > 255 {
		factor = 255
	} else if factor < 0 {
		factor = 0
	}
	color.R = uint8(math.Round(float64(color.R) * factor))
	color.G = uint8(math.Round(float64(color.G) * factor))
	color.B = uint8(math.Round(float64(color.B) * factor))
}

func getFreqPercents(fft []complex128, sampleRate ears.SampleRate, fromHz, toHz uint16) float64 {
	if len(fft) == 0 {
		return 0
	}

	idx1, idx2 := dsp.FrequencyRangeIndexes(fft, sampleRate, fromHz, toHz)
	fftCut := fft[idx1 : idx2+1]
	amp := dsp.AmplitudeMulti(fftCut)

	volume := (amp / _maxAmplitude) * 100.0
	if volume > 100 {
		volume = 100
	} else if volume < 0 {
		volume = 0
	}

	return volume
}

func frequencyRange(fromHz, toHz uint16, numBands int) []float64 {
	freqs := make([]float64, numBands)
	delta := float64((toHz - fromHz)) / float64(numBands)
	for i := 0; i < numBands; i++ {
		freqs[i] = float64(fromHz) + float64(i)*delta
	}
	return freqs
}
