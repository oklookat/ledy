package client

import (
	"encoding/binary"

	"github.com/oklookat/ledy/effect"
)

type _command uint8

const (
	// rgb array.
	_commandSetColors _command = iota
	_commandSetColorCorrection
	_commandSetColorTemperature
)

// <1 BYTE HEADER><2 BYTES UINT16 LEDS LENGTH><1 BYTE UINT8 RGB VALUES>.
//
// MAX LEDS COUNT: math.MaxUint16.
func newCommandSetColors(leds effect.LEDS) (cmd [3 + effect.LedsCount*3]uint8) {
	cmd[0] = uint8(_commandSetColors)

	ledsBytes := rgbToU8(leds)

	// <2 BYTES UINT16 LEDS LENGTH>.
	ledsLenU8 := u16toU8(uint16(len(ledsBytes)))
	cmd[1] = ledsLenU8[0]
	cmd[2] = ledsLenU8[1]

	// <1 BYTE UINT8 RGB VALUES>.
	cmdI := 3
	for i := 0; i < len(ledsBytes); i++ {
		cmd[cmdI] = ledsBytes[i]
		cmdI++
	}

	return
}

// <1 BYTE HEADER><4 BYTES UINT32 FASTLED COLOR CORRECTION VALUE>.
func newCommandSetColorCorrection(v ColorCorrection) (cmd [5]uint8) {
	return new5byteCommand(_commandSetColorCorrection, uint32(v))
}

// <1 BYTE HEADER><4 BYTES UINT32 FASTLED COLOR TEMPERATURE VALUE>.
func newCommandSetColorTemperature(v ColorTemperature) [5]uint8 {
	return new5byteCommand(_commandSetColorTemperature, uint32(v))
}

func new5byteCommand(command _command, v uint32) (cmd [5]uint8) {
	data := u32toU8(v)
	cmd[0] = uint8(command)
	cmdI := 1
	for i := 0; i < len(data); i++ {
		cmd[cmdI] = data[i]
		cmdI++
	}
	return
}

// <GRBGRBGRB>.
//
// 3 bytes per led.
func rgbToU8(leds effect.LEDS) (result [effect.LedsCount * 3]uint8) {
	// 3 bytes per led.
	resultI := 0
	for i := 0; i < effect.LedsCount; i++ {
		result[resultI] = leds[i].G
		result[resultI+1] = leds[i].R
		result[resultI+2] = leds[i].B
		resultI += 3
	}
	return
}

// 2 bytes.
func u16toU8(v uint16) (r [2]uint8) {
	binary.LittleEndian.PutUint16(r[:], v)
	return
}

// 4 bytes.
func u32toU8(v uint32) (r [4]uint8) {
	binary.LittleEndian.PutUint32(r[:], v)
	return
}
