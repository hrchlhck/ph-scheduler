package sched

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	"github.com/hrchlhck/metrics-server/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Scheduler struct {
	c    *kubernetes.Clientset
	Name string
}

func CreateScheduler(name string) *Scheduler {
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

	return &Scheduler{
		c:    clientset,
		Name: name,
	}
}

func (s *Scheduler) GetNodes() []*v1.Node {
	nl, err := s.c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	var nodeList []*v1.Node = make([]*v1.Node, 0)
	utils.CheckError(err)

	for _, node := range nl.Items {
		// taints := node.Spec.Taints

		// if len(taints) > 0 {
		// 	log.Printf("Ignoring node '%s' because of taints %+v\n", node.Name, taints)
		// 	continue
		// }
		nodeList = append(nodeList, &node)
	}

	return nodeList
}

func (s *Scheduler) GetUnscheduledPods(namespace string) []v1.Pod {
	opts := metav1.ListOptions{FieldSelector: "status.phase=Pending"}
	pods, err := s.c.CoreV1().Pods(namespace).List(context.TODO(), opts)

	var podList []v1.Pod = make([]v1.Pod, 0)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, pod := range pods.Items {
		if pod.Spec.SchedulerName == s.Name {
			podList = append(podList, pod)
		}
	}

	return podList
}

func (s *Scheduler) BestNodeForPod(pod *v1.Pod) *v1.Node {
	var schedPolicy string = string(pod.Annotations["schedulePolicy"])
	return GetNodeByPolicy(s, &schedPolicy)
}

func (s *Scheduler) bind(pod *v1.Pod, node *v1.Node) error {
	binding := v1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name: pod.Name,
		},
		Target: v1.ObjectReference{
			Kind:       "Node",
			APIVersion: "v1",
			Name:       node.Name,
		},
	}
	return s.c.CoreV1().Pods(pod.Namespace).Bind(context.TODO(), &binding, metav1.CreateOptions{})
}

func (s *Scheduler) Schedule(pod *v1.Pod) {

	node := s.BestNodeForPod(pod)
	err := s.bind(pod, node)

	if err != nil {
		log.Fatalf("Unable to schedule pod '%s' on node '%s'. Reason: %s", pod.Name, node.Name, err.Error())
	} else {
		log.Printf("Successfully scheduled pod '%s' on node '%s'.\n", pod.Name, node.Name)
	}
}
