package client

type ColorCorrection uint32

const (
	// Typical values for SMD5050 LEDs.
	ColorCorrectionTypicalLEDStrip ColorCorrection = 0xFFB0F0

	// Typical values for 8 mm "pixels on a string".
	// Also for many through-hole 'T' package LEDs.
	ColorCorrectionTypical8mmPixel ColorCorrection = 0xFFE08C

	// Uncorrected color (0xFFFFFF).
	ColorCorrectionUncorrectedColor ColorCorrection = 0xFFFFFF
)

type ColorTemperature uint32

const (
	// Black Body Radiators
	ColorTemperatureCandle         ColorTemperature = 0xFF9329 // 1900 K, 255, 147, 41
	ColorTemperatureTungsten40W    ColorTemperature = 0xFFC58F // 2600 K, 255, 197, 143
	ColorTemperatureTungsten100W   ColorTemperature = 0xFFD6AA // 2850 K, 255, 214, 170
	ColorTemperatureHalogen        ColorTemperature = 0xFFF1E0 // 3200 K, 255, 241, 224
	ColorTemperatureCarbonArc      ColorTemperature = 0xFFFAF4 // 5200 K, 255, 250, 244
	ColorTemperatureHighNoonSun    ColorTemperature = 0xFFFFFB // 5400 K, 255, 255, 251
	ColorTemperatureDirectSunlight ColorTemperature = 0xFFFFFF // 6000 K, 255, 255, 255
	ColorTemperatureOvercastSky    ColorTemperature = 0xC9E2FF // 7000 K, 201, 226, 255
	ColorTemperatureClearBlueSky   ColorTemperature = 0x409CFF // 20000 K, 64, 156, 255

	// Gaseous Light Sources
	ColorTemperatureWarmFluorescent         ColorTemperature = 0xFFF4E5 // 0 K, 255, 244, 229
	ColorTemperatureStandardFluorescent     ColorTemperature = 0xF4FFFA // 0 K, 244, 255, 250
	ColorTemperatureCoolWhiteFluorescent    ColorTemperature = 0xD4EBFF // 0 K, 212, 235, 255
	ColorTemperatureFullSpectrumFluorescent ColorTemperature = 0xFFF4F2 // 0 K, 255, 244, 242
	ColorTemperatureGrowLightFluorescent    ColorTemperature = 0xFFEFF7 // 0 K, 255, 239, 247
	ColorTemperatureBlackLightFluorescent   ColorTemperature = 0xA700FF // 0 K, 167, 0, 255
	ColorTemperatureMercuryVapor            ColorTemperature = 0xD8F7FF // 0 K, 216, 247, 255
	ColorTemperatureSodiumVapor             ColorTemperature = 0xFFD1B2 // 0 K, 255, 209, 178
	ColorTemperatureMetalHalide             ColorTemperature = 0xF2FCFF // 0 K, 242, 252, 255
	ColorTemperatureHighPressureSodium      ColorTemperature = 0xFFB74C // 0 K, 255, 183, 76

	// Uncorrected temperature (0xFFFFFF)
	ColorTemperatureUncorrectedTemperature ColorTemperature = 0xFFFFFF // 255, 255, 255
)
