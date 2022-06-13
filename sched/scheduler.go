package sched

import (
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

func WatchUnscheduledPods(scheduler *Scheduler, namespace string) {
	for {
		pods := scheduler.GetUnscheduledPods(namespace)

		if len(pods) == 0 {
			continue
		}

		for _, pod := range pods {
			mutex.Lock()
			scheduler.Schedule(&pod)
			mutex.Unlock()
		}

		time.Sleep(2 * time.Second)
	}
}
