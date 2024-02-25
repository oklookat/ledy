package dsp

import (
	"math"
	"math/cmplx"
)

var _aWeights = map[float64]float64{
	6.3:   -85.4,
	8:     -77.6,
	10:    -70.4,
	12.5:  -63.6,
	16:    -56.4,
	20:    -50.4,
	25:    -44.8,
	31.5:  -39.5,
	40:    -34.5,
	50:    -30.3,
	63:    -26.2,
	80:    -22.4,
	100:   -19.1,
	125:   -16.2,
	160:   -13.2,
	200:   -10.8,
	250:   -8.7,
	315:   -6.6,
	400:   -4.8,
	500:   -3.2,
	630:   -1.9,
	800:   -0.8,
	1000:  0.0,
	1250:  0.6,
	1600:  1.0,
	2000:  1.2,
	2500:  1.3,
	3150:  1.2,
	4000:  1.0,
	5000:  0.6,
	6300:  -0.1,
	8000:  -1.1,
	10000: -2.5,
	12500: -4.3,
	16000: -6.7,
	20000: -9.3,
}

// ApplyWeightFilter применяет весовой фильтр к сигналу FFT
func ApplyWeightFilter(fftSignal []complex128, sampleRate float64) {
	fftSize := float64(len(fftSignal))
	frequencyStep := sampleRate / fftSize
	for freq, db := range _aWeights {
		index := int(freq / frequencyStep)
		if index < len(fftSignal) {
			fftSignal[index] *= cmplx.Rect(1, dbToAmplitude(db))
		} else {
			// Найти ближайшее значение
			closestFreq := findClosestFrequency(freq, frequencyStep, len(fftSignal))
			fftSignal[closestFreq] *= cmplx.Rect(1, dbToAmplitude(db))
		}
	}
}

// dbToAmplitude конвертирует значение в децибелах в амплитуду
func dbToAmplitude(db float64) float64 {
	return math.Pow(10, db/20)
}

// findClosestFrequency находит ближайшую частоту к заданной
func findClosestFrequency(freq, frequencyStep float64, fftSize int) int {
	index := int(freq / frequencyStep)
	if index < 0 {
		return 0
	}
	if index >= fftSize {
		return fftSize - 1
	}
	return index
}
