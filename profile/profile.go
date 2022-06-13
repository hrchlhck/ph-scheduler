package profile

import (
	"log"
	"sync"
	"time"
)

type NodeProfile struct {
	Name     string
	Counter  int64
	Weights  map[string]float64
	Averages map[string]float64
}

func incorporateToMean(currentAverage, value float64, num int64) float64 {
	return currentAverage + ((value - currentAverage) / float64(num+1))
}

func CreateNode(name string, weights map[string]float64) *NodeProfile {
	return &NodeProfile{
		Name:     name,
		Counter:  0,
		Weights:  weights,
		Averages: make(map[string]float64),
	}
}

func (np *NodeProfile) Incorporate(m *Metrics) {
	if np.Counter != 0 {
		np.Averages["cpu"] = incorporateToMean(m.Cpu.LoadAvg1, m.Cpu.LoadAvg1, np.Counter)
		np.Averages["network"] = incorporateToMean(float64(m.Network.RxBytes), float64(m.Network.RxBytes), np.Counter)
		np.Averages["memory"] = incorporateToMean(float64(m.Memory.Free), float64(m.Memory.Free), np.Counter)
		np.Averages["disk"] = incorporateToMean(float64(m.Disk.ReadIO), float64(m.Disk.ReadIO), np.Counter)
		return
	}

	np.Averages["cpu"] = incorporateToMean(np.Averages["cpu"], m.Cpu.LoadAvg1, np.Counter)
	np.Averages["network"] = incorporateToMean(np.Averages["network"], float64(m.Network.RxBytes), np.Counter)
	np.Averages["memory"] = incorporateToMean(np.Averages["memory"], float64(m.Memory.Free), np.Counter)
	np.Averages["disk"] = incorporateToMean(np.Averages["disk"], float64(m.Disk.ReadIO), np.Counter)

	np.Counter++
}

func (np *NodeProfile) _score(weights *[]float64, deviceNames *[]string) float64 {
	var totalScoreDevice []float64 = []float64{}

	for _, k := range *deviceNames {
		totalScoreDevice = append(totalScoreDevice, np.SumDeviceStatistics(k))
	}

	avg, err := WeightedAverage(totalScoreDevice, weights)

	if err != nil {
		log.Fatal(err.Error())
	}

	return avg
}

func (np *NodeProfile) Score(weights []float64, deviceNames []string) float64 {
	var (
		oldScore, newScore float64
		wg                 sync.WaitGroup
	)
	wg.Add(3)

	go func() {
		defer wg.Done()
		oldScore = np._score(&weights, &deviceNames)
		log.Println(oldScore)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(np.Interval) * time.Second)
	}()

	go func() {
		defer wg.Done()
		newScore = np._score(&weights, &deviceNames)
		log.Println(newScore)
	}()

	wg.Wait()

	return newScore - oldScore
}
