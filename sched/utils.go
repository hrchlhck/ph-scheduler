package sched

import (
	"context"
	"log"

	"github.com/hrchlhck/metrics-server/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func annotateNodes(clientset *kubernetes.Clientset, annotations map[string]string) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	utils.CheckError(err)

	for _, node := range nodes.Items {
		node.SetAnnotations(annotations)
		clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})

		log.Printf("Annotated node %s with %+v\n", node.Name, annotations)
	}
}
