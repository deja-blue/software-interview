package charge

import "time"

const costPerKwH float64 = 0.2

func EstimatedPriceInCentsForCharge(startTime time.Time, duration time.Duration, drawKwHPerMin float64) float64 {
	return float64(duration.Minutes()) * costPerKwH * drawKwHPerMin
}
