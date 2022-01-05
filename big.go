package chaikin

import (
	"math/big"

	"github.com/MicahParks/go-ad"
	"github.com/MicahParks/go-ma"
)

// BigChaikin represents the state of the Chaikin Oscillator.
type BigChaikin struct {
	ad      *ad.BigAD
	short   *ma.BigEMA
	long    *ma.BigEMA
	prevBuy bool
}

// NewBig creates a new Chaikin Oscillator and returns its first point along with the corresponding Accumulation
// Distribution Line point.
func NewBig(initial [LongEMA]ad.BigInput) (cha *BigChaikin, initialResult, adLine *big.Float) {
	adLinePoints := make([]*big.Float, len(initial))
	cha = &BigChaikin{}

	cha.ad, adLine = ad.NewBig(initial[0])
	adLinePoints[0] = adLine

	for i, input := range initial[1:] {
		adLinePoints[i+1] = cha.ad.Calculate(input)
	}

	_, shortSMA := ma.NewBigSMA(adLinePoints[:ShortEMA])
	cha.short = ma.NewBigEMA(ShortEMA, shortSMA, nil)

	// Catch up the short EMA to where the long EMA will be.
	var latestShortEMA *big.Float
	for _, adLine = range adLinePoints[ShortEMA:] {
		latestShortEMA = cha.short.Calculate(adLine)
	}

	_, longSMA := ma.NewBigSMA(adLinePoints)
	cha.long = ma.NewBigEMA(LongEMA, longSMA, nil)

	initialResult = new(big.Float).Sub(latestShortEMA, longSMA)

	cha.prevBuy = initialResult.Cmp(adLine) == 1

	return cha, initialResult, adLine
}

// Calculate produces the next point on the Chaikin Oscillator given the current period's information.
func (c *BigChaikin) Calculate(next ad.BigInput) (result, adLine *big.Float, buySignal *bool) {
	adLine = c.ad.Calculate(next)
	result = new(big.Float).Sub(c.short.Calculate(adLine), c.long.Calculate(adLine))
	expected := -1
	if c.prevBuy {
		expected = 1
	}
	if result.Cmp(adLine) != expected {
		buy := !c.prevBuy
		c.prevBuy = buy
		buySignal = &buy
	}
	return result, adLine, buySignal
}
