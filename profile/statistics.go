package profile

import (
	"errors"
	"math"
	"reflect"
)

type DeviceStatistics struct {
	Mean                float64
	StandardDeviation   float64
	VarianceCoefficient float64
}

func ScoreDev(ds *DeviceStatistics) float64 {
	return (ds.Mean + ds.StandardDeviation) - (ds.Mean - ds.StandardDeviation)
}

func getDeviceStats(m interface{}) DeviceStatistics {
	v := reflect.ValueOf(m)
	values := make([]int, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)
		values[i] = int(val.Int())
	}

	mean := Mean(&values)
	std := Std(&values)
	coef := VarCoef(&values)
	return DeviceStatistics{
		Mean:                mean,
		StandardDeviation:   std,
		VarianceCoefficient: coef,
	}
}

func Mean(X *[]int) float64 {
	var total int = 0
	for _, v := range *X {
		total += v
	}
	return float64(total) / float64(len(*X))
}

func Std(X *[]int) float64 {
	var mean float64 = Mean(X)
	var variance float64 = 0.0

	for i := 0; i < len(*X); i++ {
		variance += math.Pow((mean - float64((*X)[i])), float64(2))
	}

	return math.Sqrt(variance / float64(len(*X)-1))
}

func VarCoef(X *[]int) float64 {
	var mean float64 = Mean(X)
	var std float64 = Std(X)

	if mean == 0 {
		return 0
	}

	return std / mean
}

func WeightedAverage(values *[]float64, weights *[]float64) (float64, error) {
	if len(*values) != len(*weights) {
		return math.Inf(0), errors.New("values and weights must have the same size")
	}

	var total float64 = 0
	var weightSum float64 = 0

	for i := 0; i < len(*values); i++ {
		w := (*weights)[i]
		total += (*values)[i] * w
		weightSum += w
	}

	return total / weightSum, nil
}

func Max(values *[]float64) float64 {
	var max float64 = math.Inf(-1)
	for i := 0; i < len(*values); i++ {
		v := (*values)[i]
		if v > max {
			max = v
		}
	}
	return max
}

func Min(values *[]float64) float64 {
	var min float64 = math.Inf(1)
	for i := 0; i < len(*values); i++ {
		v := (*values)[i]
		if v < min {
			min = v
		}
	}
	return min
}

func MinMax(x *float64, values *[]float64) float64 {
	return (*x - Min(values)) / (Max(values) - Min(values))
}
