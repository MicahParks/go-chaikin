package chaikin

import (
	"github.com/MicahParks/go-ad"
	"github.com/MicahParks/go-ma"
)

const (
	// ShortEMA is the number of periods in the short EMA of the Accumulation Distribution Line results. For the Chaikin
	// Oscillator.
	ShortEMA = 3

	// LongEMA is the number of periods in the long EMA of the Accumulation Distribution Line results. For the Chaikin
	// Oscillator.
	LongEMA = 10
)

// Chaikin represents the state of the Chaikin Oscillator.
type Chaikin struct {
	ad    *ad.AD
	short *ma.EMA
	long  *ma.EMA
}

// New creates a new Chaikin Oscillator and returns its first point along with the corresponding Accumulation
// Distribution Line point.
func New(initial [LongEMA]ad.Input) (chaikin Chaikin, initialResult, adLine float64) {
	adLinePoints := make([]float64, len(initial))

	chaikin.ad, adLine = ad.New(initial[0])
	adLinePoints[0] = adLine

	for i, input := range initial[1:] {
		adLinePoints[i+1] = chaikin.ad.Calculate(input)
	}

	_, shortSMA := ma.NewSMA(adLinePoints[:ShortEMA])
	chaikin.short = ma.NewEMA(ShortEMA, shortSMA, 0)

	// Catch up the short EMA to where the long EMA will be.
	var latestShortEMA float64
	for _, adLine = range adLinePoints[ShortEMA:] {
		latestShortEMA = chaikin.short.Calculate(adLine)
	}

	_, longSMA := ma.NewSMA(adLinePoints)
	chaikin.long = ma.NewEMA(LongEMA, longSMA, 0)

	return chaikin, latestShortEMA - longSMA, adLine
}

// Calculate produces the next point on the Chaikin Oscillator given the current period's information.
func (c Chaikin) Calculate(next ad.Input) (result, adLine float64) {
	adLine = c.ad.Calculate(next)
	return c.short.Calculate(adLine) - c.long.Calculate(adLine), adLine
}
