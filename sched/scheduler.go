package sched

import (
	"context"
	"flag"
	"log"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hrchlhck/metrics-server/utils"
	p "github.com/hrchlhck/ph-scheduler/profile"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Scheduler struct {
	c              *kubernetes.Clientset
	Name           string
	SchedulePolicy string
	NodeScore      map[string]float64
	wg             sync.WaitGroup
}

type NodeAnnotation struct {
	Weight float32
	Max    float64
}

var MUTEX = &sync.Mutex{}

func CreateScheduler(name, policy string, annotations map[string]string, wg sync.WaitGroup) *Scheduler {
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

	annotateNodes(clientset, annotations)

	return &Scheduler{
		c:              clientset,
		NodeScore:      make(map[string]float64),
		Name:           name,
		SchedulePolicy: policy,
		wg:             wg,
	}
}

func (s *Scheduler) GetNodes() []v1.Node {
	nl, err := s.c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	var nodeList []v1.Node = make([]v1.Node, 0)
	utils.CheckError(err)

	for _, node := range nl.Items {
		taints := node.Spec.Taints

		if len(taints) > 0 {
			log.Printf("Ignoring node '%s' because of taints %+v\n", node.Name, taints)
			continue
		}
		nodeList = append(nodeList, node)
	}

	return nodeList
}

func (s *Scheduler) GetMapNodes() map[string]v1.Node {
	nodes := s.GetNodes()
	var ret map[string]v1.Node = make(map[string]v1.Node)
	for _, node := range nodes {
		ret[node.Name] = node
	}
	return ret
}

func (s *Scheduler) GetUnscheduledPods(namespace string) []v1.Pod {
	opts := metav1.ListOptions{FieldSelector: "status.phase=Pending"}
	pods, err := s.c.CoreV1().Pods(namespace).List(context.TODO(), opts)

	var podList []v1.Pod = make([]v1.Pod, 0)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, pod := range pods.Items {
		if pod.Spec.SchedulerName == s.Name && pod.Spec.NodeName == "" {
			podList = append(podList, pod)
		}
	}

	return podList
}

func (s *Scheduler) scoreNodes(nodeProfiles *map[string]p.NodeProfile, nodes *[]v1.Node) map[string]float64 {
	var ret map[string]float64 = make(map[string]float64)

	for _, node := range *nodes {
		weights := getNodeWeights(&node)

		metrics := p.Get("http://" + node.Status.Addresses[0].Address + ":30001" + "/os/")
		np := (*nodeProfiles)[node.Name]
		np.Incorporate(metrics)

		score := np.Score(weights)

		ret[node.Name] = score
	}

	return ret
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

func (s *Scheduler) Schedule(pod *v1.Pod, node *v1.Node) {
	err := s.bind(pod, node)

	if err != nil {
		log.Fatalf("Unable to schedule pod '%s' on node '%s'. Reason: %s", pod.Name, node.Name, err.Error())
	} else {
		log.Printf("Successfully scheduled pod '%s' on node '%s'.\n", pod.Name, node.Name)
	}
}

func getNodeWeights(n *v1.Node) map[string]float64 {
	annotations := n.Annotations

	newAnnotations := make(map[string]float64)

	for k, v := range annotations {
		if strings.HasPrefix(k, "ph.weight") {
			key := strings.Split(k, "/")[1]
			newAnnotations[key] = utils.ToFloat(v, 64)
		}
	}

	return newAnnotations
}

func createNodeProfiles(nodes []v1.Node) map[string]p.NodeProfile {
	var ret map[string]p.NodeProfile = make(map[string]p.NodeProfile)
	for _, node := range nodes {
		weights := getNodeWeights(&node)
		ret[node.Name] = *p.CreateNode(node.Name, weights)
	}
	return ret
}

func (s *Scheduler) getBestNode(scores map[string]float64, policy string) v1.Node {
	nodes := s.GetMapNodes()

	log.Println(scores)

	switch policy {
	case "bestfit":
		var min float64 = 99999999999
		var retNode v1.Node
		for node, score := range scores {
			if score < min {
				retNode = nodes[node]
				min = score
			}
		}
		return retNode
	case "worstfit":
		var retNode v1.Node
		var max float64 = -999999999999

		for node, score := range scores {
			if score > max {
				retNode = nodes[node]
				max = score
			}
		}
		return retNode
	case "firstfit":
		for node := range scores {
			return nodes[node]
		}
	}

	return v1.Node{}
}

func (s *Scheduler) watchUnscheduledPods() <-chan v1.Pod {
	pods := make(chan v1.Pod)

	go func() {
		for {
			unscheduled := s.GetUnscheduledPods("default")

			for _, pod := range unscheduled {
				log.Println("Got unscheduled pod", pod.Name)
				pods <- pod
			}

			time.Sleep(1 * time.Second)
		}
	}()
	return pods
}

func MonitorUnscheduledPods(s *Scheduler) {
	pods := s.watchUnscheduledPods()

	for {
		s.wg.Wait()

		bestNode := s.getBestNode(s.NodeScore, s.SchedulePolicy)

		var pod v1.Pod = <-pods

		s.Schedule(&pod, &bestNode)

		log.Printf("Selected node '%s' (score=%.3f) based on policy '%s'\n", bestNode.Name, s.NodeScore[bestNode.Name], s.SchedulePolicy)

		time.Sleep(1 * time.Second)
	}
}

func (s *Scheduler) Start() {
	log.Printf("Starting '%s' scheduler with policy '%s'\n", s.Name, s.SchedulePolicy)

	nodes := s.GetNodes()
	nodeProfiles := createNodeProfiles(nodes)
	doneFlag := false

	for {
		oldScores := s.scoreNodes(&nodeProfiles, &nodes)
		time.Sleep(5 * time.Second)
		newScores := s.scoreNodes(&nodeProfiles, &nodes)

		if len(s.NodeScore) > 0 && !doneFlag {
			s.wg.Done()
			doneFlag = true
		}

		MUTEX.Lock()
		for k := range newScores {
			newScores[k] = math.Abs(newScores[k] - oldScores[k])
		}

		s.NodeScore = newScores

		MUTEX.Unlock()
	}
}
