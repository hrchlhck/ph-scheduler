package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	"github.com/hrchlhck/metrics-server/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	utils.CheckError(err)

	clientset, err := kubernetes.NewForConfig(config)
	utils.CheckError(err)

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	utils.CheckError(err)

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

	for _, node := range nodes.Items {
		node.SetAnnotations(annotations)
		clientset.CoreV1().Nodes().Update(context.TODO(), &node, v1.UpdateOptions{})

		log.Println(node.GetAnnotations())
	}

}
