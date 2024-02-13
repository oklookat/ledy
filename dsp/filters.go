package dsp

type ExpFilter struct {
	Value, AlphaDecay, AlphaRise float64
}

func NewExpFilter(value, alphaDecay, alphaRise float64) *ExpFilter {
	return &ExpFilter{
		Value:      value,
		AlphaDecay: alphaDecay,
		AlphaRise:  alphaDecay,
	}
}

func (e *ExpFilter) Update(newValue float64) float64 {
	alpha := e.AlphaDecay
	if newValue > e.Value {
		alpha = e.AlphaRise
	}
	e.Value = alpha*newValue + (1.0-alpha)*e.Value
	return e.Value
}
