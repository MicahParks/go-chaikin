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

// BigResult holds the results of a BigChaikin calculation.
type BigResult struct {
	ADLine      *big.Float
	BuySignal   *bool
	ChaikinLine *big.Float
}

// NewBig creates a new Chaikin Oscillator and returns its first point along with the corresponding Accumulation
// Distribution Line point.
func NewBig(initial [LongEMA]ad.BigInput) (*BigChaikin, BigResult) {
	adLinePoints := make([]*big.Float, len(initial))
	cha := &BigChaikin{}

	var adLine *big.Float
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

	result := BigResult{
		ADLine:      adLine,
		BuySignal:   nil,
		ChaikinLine: new(big.Float).Sub(latestShortEMA, longSMA),
	}

	cha.prevBuy = result.ChaikinLine.Cmp(adLine) == 1

	return cha, result
}

// Calculate produces the next point on the Chaikin Oscillator given the current period's information.
func (c *BigChaikin) Calculate(next ad.BigInput) BigResult {
	adLine := c.ad.Calculate(next)
	result := new(big.Float).Sub(c.short.Calculate(adLine), c.long.Calculate(adLine))
	expected := -1
	if c.prevBuy {
		expected = 1
	}
	var buySignal *bool
	if result.Cmp(adLine) != expected {
		buy := !c.prevBuy
		c.prevBuy = buy
		buySignal = &buy
	}
	return BigResult{
		ADLine:      adLine,
		BuySignal:   buySignal,
		ChaikinLine: result,
	}
}
