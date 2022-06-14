package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hrchlhck/ph-scheduler/sched"
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

func main() {
	sn := checkArgs()
	s := sched.CreateScheduler(sn, "bestfit")
	s.Start()
}
