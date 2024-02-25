package effect

import (
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"

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

	channel    ears.Channel
	channelFFT []complex128

	ranges         []*spectrumRange
	gradientColors []color.Color
	gradient       colorgrad.Gradient
}

func NewSpectrum(sampleRate ears.SampleRate, alphaDecay, alphaRise float64, fromHz, toHz uint16, numBands uint8) *Spectrum {
	// Gen ranges.
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

	// Gen colors.
	colors := make([]color.Color, numBands)
	for i := range colors {
		colors[i] = randomColor()
	}
	grad, err := colorgrad.NewGradient().
		Colors(
			colors...,
		).
		Build()
	if err != nil {
		return nil
	}

	return &Spectrum{
		sampleRate: sampleRate,
		alphaDecay: alphaDecay,
		alphaRise:  alphaRise,
		fromHz:     fromHz,
		toHz:       toHz,
		numBands:   numBands,

		ranges: ranges,

		gradient:       grad,
		gradientColors: colors,
	}
}

func (s *Spectrum) Visualize(data *ears.Heard, leds LEDS) {
	s.processEars(data)

	rangeSize := len(leds) / len(s.ranges)

	var wg sync.WaitGroup
	for i := range s.ranges {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			s.ranges[i].process(s.channelFFT)

			// Range idx to leds idx.
			startIdx := i * rangeSize
			endIdx := startIdx + rangeSize
			if i == len(s.ranges)-1 {
				endIdx = len(leds)
			}

			if s.ranges[i].percentsDelta > 50 {
				s.regenerateColor(i, i)
			}

			currentIdxPercents := (100 * float64(i)) / float64(len(s.ranges))
			colorAt := s.gradient.At(currentIdxPercents / 100)
			r, g, b := colorAt.RGB255()

			for x := startIdx; x < endIdx; x++ {
				if x > len(leds) {
					break
				}
				leds[x].R = r
				leds[x].G = g
				leds[x].B = b
				changeBrightness(leds[x], s.ranges[i].percents)
			}
		}(i)
	}
	wg.Wait()
}

func (s *Spectrum) regenerateColor(fromIdx, toIdx int) {
	if toIdx > len(s.gradientColors)-1 {
		return
	}
	for i := fromIdx; i < toIdx+1; i++ {
		s.gradientColors[i] = randomColor()
	}
	grad, _ := colorgrad.NewGradient().
		Colors(
			s.gradientColors...,
		).
		Build()
	s.gradient = grad
}

func (s *Spectrum) processEars(data *ears.Heard) {
	s.channel = data.Mono
	s.channelFFT = fft.FFTReal(s.channel[:len(s.channel)/2])
	dsp.ApplyWeightFilter(s.channelFFT, float64(s.sampleRate))
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomColor() color.Color {
	return color.RGBA{
		R: uint8(rnd.Intn(255)),
		G: uint8(rnd.Intn(255)),
		B: uint8(rnd.Intn(255)),
		A: 255,
	}
}

type spectrumRange struct {
	filter                *dsp.ExpFilter
	sampleRate            ears.SampleRate
	fromHz, toHz          uint16
	alphaDecay, alphaRise float64

	idx1, idx2                            int
	percents, prevPercents, percentsDirty float64
	percentsDelta                         float64
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
	s.percentsDelta = math.Abs(s.percentsDirty - float64(s.prevPercents))
	s.prevPercents = s.percentsDirty
}
