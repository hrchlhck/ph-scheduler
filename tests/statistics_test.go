package tests

import (
	hs "horchulhack-scheduler/profile"
	"math"
	"testing"
)

func TestWeightedAverageError(t *testing.T) {
	values := []float64{1, 2, 3}
	weights := []float64{1, 2, 3, 4}

	avg, err := hs.WeightedAverage(&values, &weights)

	if err != nil && avg != math.Inf(0) {
		t.Errorf("Error not empty")
	}
}

func TestWeightedAverage(t *testing.T) {
	values := []float64{1, 2, 3, 4}
	weights := []float64{0.1, 0.3, 0.5, 0.1}

	avg, _ := hs.WeightedAverage(&values, &weights)

	expected := 2.6

	if avg != expected {
		t.Errorf("Weighted average expected is %f. Got %f instead\n.", expected, avg)
	}
}
