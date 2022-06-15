package main

import (
	"log"
	"os"

	"github.com/hrchlhck/ph-scheduler/sched"
	"sync"
)

func checkArgs() string {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <scheduler name>", os.Args[0])
	}

	return os.Args[1]
}

func main() {
	var annotations map[string]string = map[string]string{
		"ph.max/cpu":        "8",
		"ph.max/memory":     "0.75",
		"ph.max/network":    "0.75",
		"ph.max/disk":       "0.75",
		"ph.weight/cpu":     "3",
		"ph.weight/memory":  "2",
		"ph.weight/network": "2",
		"ph.weight/disk":    "1",
	}
	var wg sync.WaitGroup

	// Wait for the scheduler score all nodes after N seconds
	wg.Add(1)

	sn := checkArgs()
	s := sched.CreateScheduler(sn, "bestfit", annotations, wg)

	go sched.MonitorUnscheduledPods(s)
	
	s.Start()

}
