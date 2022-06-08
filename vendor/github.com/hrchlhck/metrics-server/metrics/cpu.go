package metrics

import (
	"github.com/hrchlhck/metrics-server/utils"
)

type CPUMetrics struct {
	User             int
	Nice             int
	System           int
	Idle             int
	SoftIRQ          int
	IRQ              int
	IOWait           int
	ContextSwitches  int
	Processes        int
	AliveProcesses   int
	BlockedProcesses int
}

func GetCPUStats() interface{} {
	var cpu CPUMetrics = CPUMetrics{}

	for _, fields := range utils.GetFields("/proc/stat", false) {
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "cpu":
			cpu.User = utils.ToInt(fields[1])
			cpu.Nice = utils.ToInt(fields[2])
			cpu.System = utils.ToInt(fields[3])
			cpu.Idle = utils.ToInt(fields[4])
			cpu.SoftIRQ = utils.ToInt(fields[5])
			cpu.IRQ = utils.ToInt(fields[6])
			cpu.IOWait = utils.ToInt(fields[7])

		case "ctxt":
			cpu.ContextSwitches = utils.ToInt(fields[1])

		case "processes":
			cpu.Processes = utils.ToInt(fields[1])

		case "procs_running":
			cpu.AliveProcesses = utils.ToInt(fields[1])

		case "procs_blocked":
			cpu.BlockedProcesses = utils.ToInt(fields[1])
		}
	}

	return cpu
}
