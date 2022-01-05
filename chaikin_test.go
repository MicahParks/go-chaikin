package chaikin_test

import (
	"testing"

	"github.com/MicahParks/go-ad"

	"github.com/MicahParks/go-chaikin"
)

const (
	badBuySignal     = "The buy signal was incorrect."
	badNumericResult = "The numeric result was incorrect."
)

func BenchmarkChaikin_Calculate(b *testing.B) {
	initial := [chaikin.LongEMA]ad.Input{}
	for i := 0; i < chaikin.LongEMA; i++ {
		initial[i] = ad.Input{
			Close:  closePrices[i],
			Low:    low[i],
			High:   high[i],
			Volume: volume[i],
		}
	}

	cha, _ := chaikin.New(initial)

	for i := range open[chaikin.LongEMA:] {
		i += chaikin.LongEMA

		cha.Calculate(ad.Input{
			Close:  closePrices[i],
			Low:    low[i],
			High:   high[i],
			Volume: volume[i],
		})
	}
}

func BenchmarkBigChaikin_Calculate(b *testing.B) {
	initial := [chaikin.LongEMA]ad.BigInput{}
	for i := 0; i < chaikin.LongEMA; i++ {
		initial[i] = ad.BigInput{
			Close:  bigClose[i],
			Low:    bigLow[i],
			High:   bigHigh[i],
			Volume: bigVolume[i],
		}
	}

	cha, _ := chaikin.NewBig(initial)

	for i := range bigOpen[chaikin.LongEMA:] {
		i += chaikin.LongEMA

		cha.Calculate(ad.BigInput{
			Close:  bigClose[i],
			Low:    bigLow[i],
			High:   bigHigh[i],
			Volume: bigVolume[i],
		})
	}
}

func TestChaikin_Calculate(t *testing.T) {
	initial := [chaikin.LongEMA]ad.Input{}
	for i := 0; i < chaikin.LongEMA; i++ {
		initial[i] = ad.Input{
			Close:  closePrices[i],
			Low:    low[i],
			High:   high[i],
			Volume: volume[i],
		}
	}

	cha, result := chaikin.New(initial)
	if result.ChaikinLine != chaikinResults[0] {
		t.FailNow()
	}

	for i := range open[chaikin.LongEMA:] {
		i += chaikin.LongEMA

		result = cha.Calculate(ad.Input{
			Close:  closePrices[i],
			Low:    low[i],
			High:   high[i],
			Volume: volume[i],
		})

		if result.ChaikinLine != chaikinResults[i-chaikin.LongEMA+1] {
			t.Errorf(badNumericResult)
			t.FailNow()
		}

		expected := buySignals[i-chaikin.LongEMA]
		if result.BuySignal == nil && expected == nil {
			continue
		}
		if result.BuySignal == nil || expected == nil || *result.BuySignal != *expected {
			t.Errorf(badBuySignal)
			t.FailNow()
		}
	}
}

func TestBigChaikin_Calculate(t *testing.T) {
	initial := [chaikin.LongEMA]ad.BigInput{}
	for i := 0; i < chaikin.LongEMA; i++ {
		initial[i] = ad.BigInput{
			Close:  bigClose[i],
			Low:    bigLow[i],
			High:   bigHigh[i],
			Volume: bigVolume[i],
		}
	}

	cha, result := chaikin.NewBig(initial)
	if result.ChaikinLine.Cmp(bigChaikinResults[0]) != 0 {
		t.FailNow()
	}

	for i := range bigOpen[chaikin.LongEMA:] {
		i += chaikin.LongEMA

		result = cha.Calculate(ad.BigInput{
			Close:  bigClose[i],
			Low:    bigLow[i],
			High:   bigHigh[i],
			Volume: bigVolume[i],
		})

		if result.ChaikinLine.Cmp(bigChaikinResults[i-chaikin.LongEMA+1]) != 0 {
			t.Errorf(badNumericResult)
			t.FailNow()
		}

		expected := buySignals[i-chaikin.LongEMA]
		if result.BuySignal == nil && expected == nil {
			continue
		}
		if result.BuySignal == nil || expected == nil || *result.BuySignal != *expected {
			t.Errorf(badBuySignal)
			t.FailNow()
		}
	}
}
