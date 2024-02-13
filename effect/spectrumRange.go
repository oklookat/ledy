package effect

import (
	"math"

	"github.com/mazznoer/colorgrad"
	"github.com/mjibson/go-dsp/fft"
	"github.com/oklookat/ledy/dsp"
	"github.com/oklookat/ledy/ears"
)

type Spectrum struct {
	sampleRate            ears.SampleRate
	alphaDecay, alphaRise float64
	fromHz, toHz          uint16
	numBands              uint8

	ranges []*spectrumRange
	maxAmp float64

	channel    ears.Channel
	channelFFT []complex128
}

func NewSpectrum(sampleRate ears.SampleRate, alphaDecay, alphaRise float64, fromHz, toHz uint16, numBands uint8) *Spectrum {
	freqs := frequencyRange(fromHz, toHz, int(numBands))
	ranges := make([]*spectrumRange, 0, len(freqs))
	for i := range freqs {
		fFromHz := freqs[i]
		fToHz := 0.0
		if i+1 >= len(freqs) {
			fToHz = fFromHz
		} else {
			fToHz = freqs[i+1]
		}
		ranges = append(ranges, newSpectrumRange(sampleRate, uint16(fFromHz), uint16(fToHz), alphaDecay, alphaRise))
	}
	return &Spectrum{
		sampleRate: sampleRate,
		alphaDecay: alphaDecay,
		alphaRise:  alphaRise,
		fromHz:     fromHz,
		toHz:       toHz,
		numBands:   numBands,

		ranges: ranges,
	}
}

func (s *Spectrum) Visualize(data *ears.Heard, leds LEDS) {
	grad := colorgrad.Sinebow()

	s.processEars(data)

	rangeSize := len(leds) / len(s.ranges)

	for i := range s.ranges {
		s.ranges[i].process(s.channelFFT)

		perc := (100 * float64(i)) / float64(len(s.ranges))
		colr := grad.At(perc / 100)
		r, g, b := colr.RGB255()

		startIdx := i * rangeSize
		endIdx := startIdx + rangeSize

		if i == len(s.ranges)-1 {
			endIdx = len(leds)
		}

		for x := startIdx; x < endIdx; x++ {
			if x > len(leds) {
				break
			}
			leds[x].R = r
			leds[x].G = g
			leds[x].B = b
			changeBrightness(leds[x], s.ranges[i].percents)
		}
	}
}

func (s *Spectrum) processEars(data *ears.Heard) {
	s.channel = data.Mono
	s.channelFFT = fft.FFTReal(s.channel[:])
	for _, ff := range s.channelFFT {
		amp := dsp.AmplitudeSingle(ff)
		if amp > s.maxAmp {
			s.maxAmp = amp
		}
	}
}

type spectrumRange struct {
	filter                *dsp.ExpFilter
	sampleRate            ears.SampleRate
	fromHz, toHz          uint16
	alphaDecay, alphaRise float64

	idx1, idx2                            int
	percents, prevPercents, percentsDirty float64
	delta                                 float64
}

func newSpectrumRange(
	sampleRate ears.SampleRate,
	fromHz, toHz uint16,
	alphaDecay, alphaRise float64) *spectrumRange {
	return &spectrumRange{
		sampleRate: sampleRate,
		fromHz:     fromHz,
		toHz:       toHz,
		alphaDecay: alphaDecay,
		alphaRise:  alphaRise,
		filter:     dsp.NewExpFilter(0, alphaDecay, alphaRise),
	}
}

func (s *spectrumRange) process(fft []complex128) {
	idx1, idx2 := dsp.FrequencyRangeIndexes(fft, s.sampleRate, s.fromHz, s.toHz)
	s.idx1 = idx1
	s.idx2 = idx2

	s.percentsDirty = getFreqPercents(fft, s.sampleRate, s.fromHz, s.toHz)
	s.percents = s.filter.Update(s.percentsDirty)
	s.delta = math.Abs(s.percentsDirty - float64(s.prevPercents))
	s.prevPercents = s.percentsDirty
}
