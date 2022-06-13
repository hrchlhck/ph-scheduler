package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	s := sched.CreateScheduler(sn)
	s.Start()

	for {
		for _, node := range s.GetNodes() {
			var addr string = node.Status.Addresses[0].Address
			nodeWeight[node.Name] = getNodeWeights(node)

			metrics := p.Get("http://" + addr + "/os/")
			np := p.CreateNode(node.Name, 5, nodeWeight[node.Name])
			np.Incorporate(metrics)

			score := np.Score([]float64{1, 1, 1, 1}, []string{"cpu", "memory", "disk", "network"})

			log.Println(score)

			if score < minUsage.Score {
				minUsage = nodeTuple{node, score}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
