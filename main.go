package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"time"

	"github.com/hrchlhck/hrchlhck-scheduler/profile"
	m "github.com/hrchlhck/metrics-server/metrics"
	v1 "k8s.io/api/core/v1"
)

type nodeTuple struct {
	Node  *v1.Node
	Score float64
}

func checkArgs() string {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("Usage: %s <scheduler name>", os.Args[0]))
	}

	return os.Args[1]
}

func Subtract(c1, c2 *m.CPUMetrics) m.CPUMetrics {
	return m.CPUMetrics{
		User:             c1.User - c2.User,
		Nice:             c1.Nice - c2.Nice,
		System:           c1.System - c2.System,
		Idle:             c1.Idle - c2.Idle,
		SoftIRQ:          c1.SoftIRQ - c2.SoftIRQ,
		IRQ:              c1.IRQ - c2.IRQ,
		IOWait:           c1.IOWait - c2.IOWait,
		ContextSwitches:  c1.ContextSwitches - c2.ContextSwitches,
		Processes:        c1.Processes - c2.Processes,
		AliveProcesses:   c1.AliveProcesses - c2.AliveProcesses,
		BlockedProcesses: c1.BlockedProcesses - c2.BlockedProcesses,
	}
}

func PercDiff(a, b int) float64 {
	var v1, v2 float64 = float64(a), float64(b)
	var diff float64 = math.Abs(v1 - v2)

	if diff == 0 {
		return 0
	}

	return diff / ((v1 + v2) / 2) * 100
}

func Rate(c1, c2 m.CPUMetrics) float64 {
	m1 := reflect.ValueOf(c1)
	m2 := reflect.ValueOf(c2)
	var sum float64 = 0

	for i := 0; i < m1.NumField(); i++ {
		val1 := m1.Field(i)
		val2 := m2.Field(i)
		sum += PercDiff(int(val1.Int()), int(val2.Int()))
	}

	return sum / float64(m1.NumField())
}

func main() {
	// sn := checkArgs()
	// s := sched.CreateScheduler(sn)
	// addr := "http://" + s.GetNodes()[0].Status.Addresses[0].Address + "/os/"
	addr := "http://172.17.0.2/os/"

	for {
		metrics := profile.Get(addr).Cpu

		time.Sleep(5 * time.Second)

		metrics1 := profile.Get(addr).Cpu

		fmt.Println(Rate(metrics, metrics1))
	}

	// log.Printf("Starting %s scheduler\n", sn)
	// sched.WatchUnscheduledPods(s, "default")
	// var minUsage nodeTuple = nodeTuple{nil, math.Inf(99999)}
	// for {
	// 	for _, node := range s.GetNodes() {
	// 		// var addr string = node.Status.Addresses[0].Address
	// 		addr := "172.17.0.2"

	// 		metrics := *hs.Get("http://" + addr + "/os/")
	// 		np := hs.CreateNode(node.Name, 10, 5)
	// 		np.Incorporate(metrics.Statistics())
	// 		score := np.Score([]float64{1, 1, 1, 1}, []string{"cpu", "memory", "disk", "network"})
	// 		// np._score()
	// 		log.Println(score)

	// 		if score < minUsage.Score {
	// 			minUsage = nodeTuple{node, score}
	// 		}
	// 	}
	// 	time.Sleep(2 * time.Second)
	// }
}
