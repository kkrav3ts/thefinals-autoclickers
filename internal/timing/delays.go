package timing

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

// GenerateClickDelays returns count delays drawn from a normal distribution
// and clamped to [minVal, maxVal]. Values are milliseconds.
func GenerateClickDelays(count int, mean, stdDev, minVal, maxVal float64) ([]time.Duration, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}
	if stdDev < 0 {
		return nil, fmt.Errorf("stddev must be zero or positive")
	}
	if minVal > maxVal {
		return nil, fmt.Errorf("min %.2fms is greater than max %.2fms", minVal, maxVal)
	}

	delays := make([]time.Duration, count)
	for i := 0; i < count; i++ {
		val := mean + stdDev*rand.NormFloat64()
		val = math.Max(minVal, math.Min(maxVal, val))
		delays[i] = time.Duration(math.Round(val)) * time.Millisecond
	}
	return delays, nil
}
