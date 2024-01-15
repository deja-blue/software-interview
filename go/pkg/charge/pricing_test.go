package charge

import (
	"testing"
	"time"
)

func TestEstimatedPriceInCentsForCharge(t *testing.T) {
	startTime := time.Now()
	duration := 30 * time.Minute
	drawKwHPerMin := 0.5

	expectedPrice := float64(duration.Minutes()) * costPerKwH * drawKwHPerMin
	actualPrice := EstimatedPriceInCentsForCharge(startTime, duration, drawKwHPerMin)

	if actualPrice != expectedPrice {
		t.Errorf("Expected price: %f, but got: %f", expectedPrice, actualPrice)
	}
}
