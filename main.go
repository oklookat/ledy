package main

import (
	"os"
	"time"

	"github.com/oklookat/ledy/client"
	"github.com/oklookat/ledy/ears"
	"github.com/oklookat/ledy/effect"
)

func main() {
	cl := client.New()
	chk(cl.Connect())
	println("Connected.")
	cl.SetColorCorrection(client.ColorCorrectionTypicalLEDStrip)
	cl.SetColorTemperature(client.ColorTemperatureCoolWhiteFluorescent)

	leds := effect.NewLEDS()
	var ledSpec *effect.Spectrum

	listener := ears.New(func(sampleRate ears.SampleRate) {
		ledSpec = effect.NewSpectrum(sampleRate, 0.3, 0.3, 30, 1000, 14)
	})

	chk(listener.Listen(func(h *ears.Heard) {
		ledSpec.Visualize(h, leds)
		_ = cl.SetColors(leds)
	}))

	for {
		time.Sleep(24 * time.Hour)
	}
}

func chk(err error) {
	if err == nil {
		return
	}
	println(err.Error())
	os.Exit(1)
}
