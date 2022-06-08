package metrics

import "github.com/hrchlhck/metrics-server/utils"

type NetworkMetrics struct {
	RxBytes      int
	RxPackets    int
	RxErrors     int
	RxDrop       int
	RxFifo       int
	RxFrame      int
	RxCompressed int
	RxMulticast  int
	TxBytes      int
	TxPackets    int
	TxErrors     int
	TxDrop       int
	TxFifo       int
	TxFrame      int
	TxCompressed int
	TxMulticast  int
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

func GetNetworkStats() interface{} {
	var net NetworkMetrics = NetworkMetrics{}
	var data []string = getIface("eno1", utils.GetFields("/proc/net/dev", false))

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
	}

	return net
}
