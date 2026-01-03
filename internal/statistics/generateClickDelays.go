package statistics

import (
	"math"
	"math/rand/v2"
	"time"
)

// GenerateClickDelays generates a list of realistic click delays.
// The delays follow a normal distribution clustered around 60-85ms, with bounds of 50-105ms.
//
// Parameters:
//   - count  : number of delays to generate
//   - mean   : midpoint of delay cluster
//   - stdDev : standard deviation to create the delay cluster
//   - minVal : minimum delay
//   - maxVal : maximum delay
//
// Returns: slice of time.Duration representing delays (already multiplied by time.Millisecond)
func GenerateClickDelays(count int, mean float64, stdDev float64, minVal float64, maxVal float64) []time.Duration {
	if count <= 0 {
		return []time.Duration{}
	}

	delays := make([]time.Duration, count)

	for i := 0; i < count; i++ {
		// Generate value from normal distribution
		val := mean + stdDev*rand.NormFloat64()

		// Clamp to bounds [minVal, maxVal]
		val = math.Max(minVal, math.Min(maxVal, val))

		// Convert to time.Duration by multiplying milliseconds
		delays[i] = time.Duration(math.Round(val)) * time.Millisecond
	}

	return delays
}
