package main

import (
	"fmt"
	"horchulhack-scheduler/sched"
	"log"
	"os"

	"github.com/hrchlhck/metrics-server/utils"
)

func checkArgs() (string, int, int) {
	if len(os.Args) < 4 {
		log.Fatal(fmt.Sprintf("Usage: %s <metrics server ip> <window size> <interval>", os.Args[0]))
	}

	return os.Args[1], utils.ToInt(os.Args[2]), utils.ToInt(os.Args[3])
}

func main() {
	s := sched.CreateScheduler("horchulhack")

	pods := s.GetUnscheduledPods("default")
	nodes, err := s.GetNodes()
	utils.CheckError(err)

	for _, pod := range pods {
		s.BestNodeForPod(&pod)
		s.Schedule(&pod, &nodes.Items[0])
	}
}
