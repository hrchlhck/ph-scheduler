package metrics

import (
	"github.com/hrchlhck/metrics-server/utils"
)

type MemoryMetrics struct {
	Free              int
	Available         int
	Buffers           int
	Cached            int
	SwapCached        int
	ActivePages       int
	InactivePages     int
	ActiveAnonPages   int
	InactiveAnonPages int
	Mapped            int
	// KernelStack       int
}

func GetMemoryStats() interface{} {
	var mem MemoryMetrics = MemoryMetrics{}
	var data [][]string = utils.GetFields("/proc/meminfo", false)

	if len(data) != 0 {
		mem.Free = utils.ToInt(data[1][1])
		mem.Available = utils.ToInt(data[2][1])
		mem.Buffers = utils.ToInt(data[3][1])
		mem.Cached = utils.ToInt(data[4][1])
		mem.SwapCached = utils.ToInt(data[5][1])
		mem.ActivePages = utils.ToInt(data[6][1])
		mem.InactivePages = utils.ToInt(data[7][1])
		mem.ActiveAnonPages = utils.ToInt(data[8][1])
		mem.InactiveAnonPages = utils.ToInt(data[9][1])
		mem.Mapped = utils.ToInt(data[19][1])
		return mem
	}
	return nil
}
