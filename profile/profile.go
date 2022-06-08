package profile

import (
	"log"
	"sync"
	"time"
)

type NodeProfile struct {
	Name       string
	Interval   int
	Counter    int
	WindowSize int
	Window     []map[string]DeviceStatistics
}

func CreateNode(name string, windowSize, interval int) *NodeProfile {
	return &NodeProfile{
		Name:       name,
		Counter:    0,
		WindowSize: windowSize,
		Window:     make([]map[string]DeviceStatistics, windowSize),
	}
}

func (np *NodeProfile) Incorporate(stats map[string]DeviceStatistics) {
	if np.Counter == np.WindowSize-1 {
		np.Decrement()

		np.Window = append(np.Window, stats)

		return
	}
	np.Window[np.Counter] = stats
	np._increment()
}

func (np *NodeProfile) Decrement() {
	// Pop first item
	_, np.Window = np.Window[0], np.Window[1:]
}

func (np *NodeProfile) _increment() {
	if np.Counter < np.WindowSize {
		np.Counter++
	}
}

func (np *NodeProfile) SumDeviceStatistics(device string) float64 {
	var total float64 = 0

	for i := 0; i < np.WindowSize-1; i++ {
		if np.Window[i] == nil {
			continue
		}

		ds := np.Window[i][device]
		total += ScoreDev(&ds)
	}

	return total
}

func (np *NodeProfile) _score(weights *[]float64, deviceNames *[]string) float64 {
	var totalScoreDevice []float64 = []float64{}

	for _, k := range *deviceNames {
		totalScoreDevice = append(totalScoreDevice, np.SumDeviceStatistics(k))
	}

	avg, err := WeightedAverage(&totalScoreDevice, weights)

	if err != nil {
		log.Fatal(err.Error())
	}

	// return MinMax(&avg, &totalScoreDevice)
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
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(np.Interval) * time.Second)
	}()

	go func() {
		defer wg.Done()
		newScore = np._score(&weights, &deviceNames)
	}()

	wg.Wait()

	return newScore - oldScore
}
