package profile

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	m "github.com/hrchlhck/metrics-server/metrics"
	utils "github.com/hrchlhck/metrics-server/utils"
)

type Metrics struct {
	Cpu     m.CPUMetrics
	Memory  m.MemoryMetrics
	Disk    m.DiskMetrics
	Network m.NetworkMetrics
}

func (m *Metrics) Statistics() map[string]DeviceStatistics {
	return map[string]DeviceStatistics{
		"cpu":     getDeviceStats(m.Cpu),
		"memory":  getDeviceStats(m.Memory),
		"disk":    getDeviceStats(m.Disk),
		"network": getDeviceStats(m.Network),
	}
}

func ReadJSON(addr, device string) []byte {
	response, err := http.Get(addr + device)

	utils.CheckError(err)

	responseData, err := ioutil.ReadAll(response.Body)
	utils.CheckError(err)

	return responseData
}

func Get(addr string) *Metrics {
	var metrics Metrics = Metrics{
		Cpu:     m.CPUMetrics{},
		Memory:  m.MemoryMetrics{},
		Disk:    m.DiskMetrics{},
		Network: m.NetworkMetrics{},
	}

	var cpu m.CPUMetrics
	var memory m.MemoryMetrics
	var disk m.DiskMetrics
	var network m.NetworkMetrics

	var metricBytes map[string][]byte = map[string][]byte{
		"cpu":     ReadJSON(addr, "cpu"),
		"memory":  ReadJSON(addr, "memory"),
		"disk":    ReadJSON(addr, "disk"),
		"network": ReadJSON(addr, "network"),
	}

	json.Unmarshal(metricBytes["cpu"], &cpu)
	json.Unmarshal(metricBytes["memory"], &memory)
	json.Unmarshal(metricBytes["disk"], &disk)
	json.Unmarshal(metricBytes["network"], &network)

	metrics.Cpu = cpu
	metrics.Memory = memory
	metrics.Disk = disk
	metrics.Network = network

	return &metrics
}
