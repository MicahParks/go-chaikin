package chaikin

import (
	"math/big"

	"github.com/MicahParks/go-ad"
	"github.com/MicahParks/go-ma"
)

// BigChaikin represents the state of the Chaikin Oscillator.
type BigChaikin struct {
	ad    *ad.BigAD
	short *ma.BigEMA
	long  *ma.BigEMA
}

// NewBig creates a new Chaikin Oscillator and returns its first point along with the corresponding Accumulation
// Distribution Line point.
func NewBig(initial [LongEMA]ad.BigInput) (chaikin BigChaikin, initialResult, adLine *big.Float) {
	adLinePoints := make([]*big.Float, len(initial))

	chaikin.ad, adLine = ad.NewBig(initial[0])
	adLinePoints[0] = adLine

	for i, input := range initial[1:] {
		adLinePoints[i+1] = chaikin.ad.Calculate(input)
	}

	_, shortSMA := ma.NewBigSMA(adLinePoints[:ShortEMA])
	chaikin.short = ma.NewBigEMA(ShortEMA, shortSMA, nil)

	// Catch up the short EMA to where the long EMA will be.
	var latestShortEMA *big.Float
	for _, adLine = range adLinePoints[ShortEMA:] {
		latestShortEMA = chaikin.short.Calculate(adLine)
	}

	_, longSMA := ma.NewBigSMA(adLinePoints)
	chaikin.long = ma.NewBigEMA(LongEMA, longSMA, nil)

	return chaikin, new(big.Float).Sub(latestShortEMA, longSMA), adLine
}

// Calculate produces the next point on the Chaikin Oscillator given the current period's information.
func (c BigChaikin) Calculate(next ad.BigInput) (result, adLine *big.Float) {
	adLine = c.ad.Calculate(next)
	return new(big.Float).Sub(c.short.Calculate(adLine), c.long.Calculate(adLine)), adLine
}
