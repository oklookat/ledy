package ears

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/gen2brain/malgo"
)

type (
	SampleRate uint32
	Channel    [BufferSize]float64
)

const BufferSize = 1024

type Heard struct {
	SampleRate SampleRate
	Mono       Channel
}

type Ears struct {
	onGotDevice func(sampleRate SampleRate)
	device      *malgo.Device
	ctx         *malgo.AllocatedContext
}

func New(onGotDevice func(sampleRate SampleRate)) *Ears {
	return &Ears{
		onGotDevice: onGotDevice,
	}
}

func (e *Ears) Listen(onData func(*Heard)) error {
	// Init.
	ctx, err := malgo.InitContext([]malgo.Backend{
		malgo.BackendWasapi,
	}, malgo.ContextConfig{}, nil)
	if err != nil {
		return err
	}
	e.ctx = ctx

	// Get default playback device.
	infos, err := ctx.Devices(malgo.Playback)
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return errors.New("no playback devices")
	}
	full, err := ctx.DeviceInfo(malgo.Playback, infos[0].ID, malgo.Shared)
	if err != nil {
		return err
	}
	if len(full.Formats) == 0 {
		return errors.New("playback device without formats")
	}

	// Setup config.
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Loopback)
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = full.Formats[0].SampleRate
	deviceConfig.PeriodSizeInFrames = BufferSize

	if e.onGotDevice != nil {
		e.onGotDevice(SampleRate(deviceConfig.SampleRate))
	}

	// Start loopback.
	dev, err := malgo.InitDevice(ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: func(pOutputSample, pInputSamples []byte, framecount uint32) {
			if onData == nil {
				return
			}
			onData(&Heard{
				SampleRate: SampleRate(deviceConfig.SampleRate),
				Mono:       convertBuffer(pInputSamples),
			})
		},
	})
	if err != nil {
		return err
	}
	e.device = dev

	return e.device.Start()
}

func (e *Ears) Unlisten() {
	if e.device != nil {
		_ = e.device.Stop()
		e.device = nil
	}
	if e.ctx != nil {
		_ = e.ctx.Uninit()
		e.ctx.Free()
		e.ctx = nil
	}
}

func convertBuffer(bytes []byte) (result [BufferSize]float64) {
	buffI := 0
	for i := 0; i < len(bytes); i += 4 {
		bits := binary.LittleEndian.Uint32(bytes[i : i+4])
		floatd := math.Float32frombits(bits)
		result[buffI] = float64(floatd)
		buffI++
	}
	return result
}
