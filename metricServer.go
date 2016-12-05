package main

import (
	"log"
	"sync"
)

type metricServer struct {
	entities int
	mutex    sync.Mutex
}

func newMetricServer(numberOfServers int) *metricServer {
	return &metricServer{
		entities: 0,
		mutex:    sync.Mutex{},
	}
}

func (ms *metricServer) flushMetrics() {
	ms.mutex.Lock()
	log.Println("Global Entities: ", ms.entities)
	ms.entities = 0
	ms.mutex.Unlock()
}

func (ms *metricServer) addMetrics(entityNumber int) {
	ms.mutex.Lock()
	ms.entities = ms.entities + entityNumber
	ms.mutex.Unlock()
}
