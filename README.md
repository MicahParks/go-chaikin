[![Go Reference](https://pkg.go.dev/badge/github.com/MicahParks/go-chaikin.svg)](https://pkg.go.dev/github.com/MicahParks/go-chaikin) [![Go Report Card](https://goreportcard.com/badge/github.com/MicahParks/go-chaikin)](https://goreportcard.com/report/github.com/MicahParks/go-chaikin)
# go-chaikin
The Chaikin Oscillator technical analysis algorithm implemented in Golang.

```go
import "github.com/MicahParks/go-chaikin"
```

## Usage
For full examples, please see the `examples` directory.

### Step 1
Gather the initial input. This is 10 periods of inputs for the Accumulation Distribution Line. It is the minimum number
of periods required to produce one result for the Chaikin Oscillator. The value, `10` is stored in `chaikin.LongEMA`
constant.

The input must be put into an array, not a slice. Make sure to fill all 10 indices of the array.
```go
initial := [chaikin.LongEMA]ad.Input{}
for i := 0; i < chaikin.LongEMA; i++ {
	initial[i] = ad.Input{
		Close:  closePrices[i],
		Low:    low[i],
		High:   high[i],
		Volume: volume[i],
	}
}
```

### Step 2
Create the Chaikin Oscillator data structure from the initial input array. This will also produce the first Chaikin
Oscillator point as well as its corresponding Accumulation Distribution Line point.
```go
cha, firstResult, adLine := chaikin.New(initial)
```

### Step 3
Use the subsequent periods to calculate the next points for the Chaikin Oscillator and Accumulation Distribution Line.
```go
result, adLine = cha.Calculate(ad.Input{
	Close:  closePrices[i],
	Low:    low[i],
	High:   high[i],
	Volume: volume[i],
})
```

## Somewhat complete example (without data)
```go
package main

import (
	"log"
	"os"

	"github.com/MicahParks/go-ad"

	"github.com/MicahParks/go-chaikin"
)

func main() {
	// Create a logger.
	logger := log.New(os.Stdout, "", 0)

	// Create the initial input.
	initial := [chaikin.LongEMA]ad.Input{}
	for i := 0; i < chaikin.LongEMA; i++ {
		initial[i] = ad.Input{
			Close:  closePrices[i],
			Low:    low[i],
			High:   high[i],
			Volume: volume[i],
		}
	}

	// Create the Chaikin Oscillator data structure as well as its first data point and the corresponding Accumulation
	// Distribution Line point.
	cha, result, adLine := chaikin.New(initial)
	logger.Printf("%.4f, %.4f", adLine, result)

	// Use every subsequent period's data to calculate the next points on the Chaikin Oscillator and Accumulation
	// Distribution Line.
	for i := range open[chaikin.LongEMA:] {
		i += chaikin.LongEMA

		result, adLine = cha.Calculate(ad.Input{
			Close:  closePrices[i],
			Low:    low[i],
			High:   high[i],
			Volume: volume[i],
		})
		logger.Printf("%.4f, %.4f", adLine, result)
	}
}
```

## Testing
There is 100% test coverage and benchmarks for this project. Here is an example benchmark result:
```
$ go test -bench .
goos: linux
goarch: amd64
pkg: github.com/MicahParks/go-chaikin
cpu: Intel(R) Core(TM) i5-9600K CPU @ 3.70GHz
BenchmarkChaikin_Calculate-6            1000000000               0.0000017 ns/op
BenchmarkBigChaikin_Calculate-6         1000000000               0.0000891 ns/op
PASS
ok      github.com/MicahParks/go-chaikin        0.004s
```

## Resources
I built and tested this package using these resources:
* [Investopedia](https://www.investopedia.com/terms/c/chaikinoscillator.asp)
* [Invest Excel](https://investexcel.net/chaikin-oscillator-spreadsheet-vba/)
