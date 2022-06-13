package metrics

import "github.com/hrchlhck/metrics-server/utils"

type DiskMetrics struct {
	ReadIO         int
	ReadMerges     int
	ReadSectors    int
	ReadTicks      int
	WriteIO        int
	WriteMerges    int
	WriteSectors   int
	WriteTicks     int
	InFlight       int
	IOTicks        int
	TimeInQueue    int
	DiscardIO      int
	DiscardMerges  int
	DiscardSectors int
	DiscardTicks   int
	FlushIO        int
	FlushTicks     int
}

func GetDiskStats() *DiskMetrics {
	var disk DiskMetrics = DiskMetrics{}
	var data []string = utils.GetFields("/sys/block/sda/stat", true)[0]

	disk.ReadIO = utils.ToInt(data[0])
	disk.ReadMerges = utils.ToInt(data[1])
	disk.ReadSectors = utils.ToInt(data[2])
	disk.ReadTicks = utils.ToInt(data[3])
	disk.WriteIO = utils.ToInt(data[4])
	disk.WriteMerges = utils.ToInt(data[5])
	disk.WriteSectors = utils.ToInt(data[6])
	disk.WriteTicks = utils.ToInt(data[7])
	disk.InFlight = utils.ToInt(data[8])
	disk.IOTicks = utils.ToInt(data[9])
	disk.TimeInQueue = utils.ToInt(data[10])
	disk.DiscardIO = utils.ToInt(data[11])
	disk.DiscardMerges = utils.ToInt(data[12])
	disk.DiscardSectors = utils.ToInt(data[13])
	disk.DiscardTicks = utils.ToInt(data[14])
	disk.FlushIO = utils.ToInt(data[15])
	disk.FlushTicks = utils.ToInt(data[16])

	return &disk
}
