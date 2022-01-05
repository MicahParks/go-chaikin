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
	ad      *ad.AD
	short   *ma.EMA
	long    *ma.EMA
	prevBuy bool
}

// Result holds the results of a Chaikin calculation.
type Result struct {
	ADLine      float64
	BuySignal   *bool
	ChaikinLine float64
}

// New creates a new Chaikin Oscillator and returns its first point along with the corresponding Accumulation
// Distribution Line point.
func New(initial [LongEMA]ad.Input) (*Chaikin, Result) {
	adLinePoints := make([]float64, len(initial))
	cha := &Chaikin{}

	var adLine float64
	cha.ad, adLine = ad.New(initial[0])
	adLinePoints[0] = adLine

	for i, input := range initial[1:] {
		adLinePoints[i+1] = cha.ad.Calculate(input)
	}

	_, shortSMA := ma.NewSMA(adLinePoints[:ShortEMA])
	cha.short = ma.NewEMA(ShortEMA, shortSMA, 0)

	// Catch up the short EMA to where the long EMA will be.
	var latestShortEMA float64
	for _, adLine = range adLinePoints[ShortEMA:] {
		latestShortEMA = cha.short.Calculate(adLine)
	}

	_, longSMA := ma.NewSMA(adLinePoints)
	cha.long = ma.NewEMA(LongEMA, longSMA, 0)

	result := Result{
		ADLine:      adLine,
		BuySignal:   nil,
		ChaikinLine: latestShortEMA - longSMA,
	}

	cha.prevBuy = result.ChaikinLine > adLine

	return cha, result
}

// Calculate produces the next point on the Chaikin Oscillator given the current period's information.
func (c *Chaikin) Calculate(next ad.Input) Result {
	adLine := c.ad.Calculate(next)
	result := c.short.Calculate(adLine) - c.long.Calculate(adLine)
	var buySignal *bool
	if result > adLine != c.prevBuy {
		buy := !c.prevBuy
		c.prevBuy = buy
		buySignal = &buy
	}
	return Result{
		ADLine:      adLine,
		BuySignal:   buySignal,
		ChaikinLine: result,
	}
}
