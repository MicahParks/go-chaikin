package chaikin

import (
	"github.com/MicahParks/go-ad"
	"github.com/MicahParks/go-ma"
)

const (
	// TODO
	ShortEMA = 3
	LongEMA  = 10
)

type Chaikin struct {
	ad    *ad.AD
	short *ma.EMA
	long  *ma.EMA
}

func New(initial [LongEMA]ad.Input) (chaikin Chaikin, initialResult, adLine float64) {
	adLinePoints := make([]float64, len(initial))
	chaikin = Chaikin{}

	chaikin.ad, adLine = ad.New(initial[0])
	adLinePoints[0] = adLine

	for i, input := range initial[1:] {
		adLinePoints[i+1] = chaikin.ad.Calculate(input)
	}

	_, shortSMA := ma.NewSMA(adLinePoints[:ShortEMA])
	short := ma.NewEMA(ShortEMA, shortSMA, 0)

	// Catch up the short EMA to where the long EMA will be.
	var latestShortEMA float64
	for _, adLine = range adLinePoints[ShortEMA:] {
		latestShortEMA = short.Calculate(adLine)
	}

	_, longSMA := ma.NewSMA(adLinePoints)
	long := ma.NewEMA(LongEMA, longSMA, 0)

	chaikin.short = short
	chaikin.long = long

	return chaikin, latestShortEMA - longSMA, adLine
}

func (c Chaikin) Calculate(next ad.Input) (result, adPoint float64) {
	adPoint = c.ad.Calculate(next)
	return c.short.Calculate(adPoint) - c.long.Calculate(adPoint), adPoint
}
