package metrics

import (
	"strings"

	"github.com/hrchlhck/metrics-server/utils"
)

type NetworkMetrics struct {
	RxBytes      int
	RxPackets    int
	RxErrors     int
	RxDrop       int
	RxFifo       int
	RxFrame      int
	RxCompressed int
	RxMulticast  int
	RxBytesRate  float64
	TxBytes      int
	TxPackets    int
	TxErrors     int
	TxDrop       int
	TxFifo       int
	TxFrame      int
	TxCompressed int
	TxMulticast  int
	TxBytesRate  float64
	Speed        int
}

func getIface(iface string, data [][]string) []string {
	for _, line := range data {
		if len(line) == 0 {
			continue
		}

		if line[0] == iface+":" {
			return line[1:]
		}
	}
	return []string{}
}

func ListIfaces(data [][]string) []string {
	var ret []string = make([]string, 0)
	for i := 2; i < len(data); i++ {
		var line []string = data[i]

		if len(line) == 0 {
			continue
		}

		newline := strings.Replace(line[0], ":", "", 1)
		ret = append(ret, newline)
	}
	return ret
}

func getIfaceSpeed(iface string) int {
	return utils.ToInt(utils.GetFields("/sys/class/net/"+iface+"/speed", true)[0][0])
}

func GetNetworkStats(iface string) *NetworkMetrics {
	var net NetworkMetrics = NetworkMetrics{}
	d := utils.GetFields("/procfs/net/dev", false)
	var data []string = getIface(iface, d)
	var speed int = getIfaceSpeed(iface)

	if len(data) != 0 {
		net.RxBytes = utils.ToInt(data[0])
		net.RxPackets = utils.ToInt(data[1])
		net.RxErrors = utils.ToInt(data[2])
		net.RxDrop = utils.ToInt(data[3])
		net.RxFifo = utils.ToInt(data[4])
		net.RxFrame = utils.ToInt(data[5])
		net.RxCompressed = utils.ToInt(data[6])
		net.RxMulticast = utils.ToInt(data[7])
		net.TxBytes = utils.ToInt(data[8])
		net.TxPackets = utils.ToInt(data[9])
		net.TxErrors = utils.ToInt(data[10])
		net.TxDrop = utils.ToInt(data[11])
		net.TxFifo = utils.ToInt(data[12])
		net.TxFrame = utils.ToInt(data[13])
		net.TxCompressed = utils.ToInt(data[14])
		net.TxMulticast = utils.ToInt(data[15])
		net.Speed = speed
	}

	return &net
}
