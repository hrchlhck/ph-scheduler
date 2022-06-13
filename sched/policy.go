package sched

import (
	"log"

	v1 "k8s.io/api/core/v1"
)

type nodeTuple struct {
	Node  *v1.Node
	Score float64
}

func bestFit(s *Scheduler) *v1.Node {
	nodes := s.GetNodes()

	// var minUsage nodeTuple = nodeTuple{nil, math.Inf(99999)}
	// for _, node := range nodes.Items {
	// 	var addr string = node.Status.Addresses[0].Address

	// 	metrics := *hs.Get("http://" + addr + "/os/")
	// 	np := hs.CreateNode(node.Name, 10, 5)
	// 	np.Incorporate(metrics.Statistics())
	// 	score := np.Score([]float64{1, 1, 1, 1}, []string{"cpu", "memory", "disk", "network"})

	// 	if score < minUsage.Score {
	// 		minUsage = nodeTuple{&node, score}
	// 	}
	// }

	for _, node := range nodes {
		log.Println(node.Name)
	}
	return nodes[0]
}

func GetNodeByPolicy(s *Scheduler, policy *string) *v1.Node {
	var ret *v1.Node

	switch *policy {
	case "bestfit":
		ret = bestFit(s)
	case "worstfit":
		ret = nil
	case "firstfit":
		ret = nil
	}

	return ret
}