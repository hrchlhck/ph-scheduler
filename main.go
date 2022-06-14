package main

import (
	"log"
	"os"

	"github.com/hrchlhck/ph-scheduler/sched"
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
		"ph.weight/cpu":     "2",
		"ph.weight/memory":  "1",
		"ph.weight/network": "2",
		"ph.weight/disk":    "3",
	}

	sn := checkArgs()
	s := sched.CreateScheduler(sn, "bestfit", annotations)

	s.Start()

	go sched.MonitorUnscheduledPods(s)

}
