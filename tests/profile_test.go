package tests

import (
	"reflect"
	"testing"
	"time"

	hs "github.com/hrchlhck/hrchlhck-scheduler/profile"

	"github.com/google/go-cmp/cmp"
)

func TestWindowIncorporateOne(t *testing.T) {
	var np *hs.NodeProfile = hs.CreateNode("node1", 10, 5)

	metrics := *hs.Get("http://172.17.0.2/os/")
	stats := metrics.Statistics()

	var expected []map[string]hs.DeviceStatistics = []map[string]hs.DeviceStatistics{
		stats,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
	np.Incorporate(stats)

	if !reflect.DeepEqual(expected, np.Window) {
		t.Errorf("Expected %v, got %v.\n", expected, np.Window)
	}
}

func TestWindowIncorporateN(t *testing.T) {
	var np *hs.NodeProfile = hs.CreateNode("node1", 10, 5)

	metrics := *hs.Get("http://172.17.0.2/os/")
	stats := metrics.Statistics()

	time.Sleep(1 * time.Second)

	stats2 := metrics.Statistics()

	var expected []map[string]hs.DeviceStatistics = []map[string]hs.DeviceStatistics{
		stats,
		stats2,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
	np.Incorporate(stats)
	np.Incorporate(stats2)

	if !reflect.DeepEqual(expected, np.Window) {
		t.Errorf("Expected %v, got %v.\n", expected, np.Window)
	}
}

func TestWindowIncorporateShift(t *testing.T) {
	var np *hs.NodeProfile = hs.CreateNode("node1", 10, 5)

	metrics := *hs.Get("http://172.17.0.2/os/")
	stats := metrics.Statistics()

	var expected []map[string]hs.DeviceStatistics = []map[string]hs.DeviceStatistics{
		metrics.Statistics(),
		metrics.Statistics(),
		metrics.Statistics(),
		metrics.Statistics(),
		metrics.Statistics(),
		metrics.Statistics(),
		metrics.Statistics(),
		metrics.Statistics(),
		metrics.Statistics(),
		metrics.Statistics(),
	}
	np.Counter = 9
	np.Window = expected
	first := np.Window[1]
	last := np.Window[len(np.Window)-1]

	np.Incorporate(stats)

	if !cmp.Equal(first, np.Window[0]) {
		t.Errorf("First element must be %v. \nGot %v instead.\n", np.Window[0], first)
	}

	if !cmp.Equal(stats, last) {
		t.Errorf("Last element must be %v. \nGot %v instead.\n", stats, last)
	}

}
