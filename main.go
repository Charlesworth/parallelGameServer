package main

import (
	"fmt"
	"math/rand"
	"time"
)

var passedChanBufSize = 10
var outOfBoundEntities chan (Entity)

func init() {
	outOfBoundEntities = make(chan Entity, passedChanBufSize)
}

type PositionServerSupervisor struct {
	positionServers []*PositionServer
}

func (pss *PositionServerSupervisor) relocateOutOfBoundsEntities() {
	for entity := range outOfBoundEntities {
		for _, ps := range pss.positionServers {
			if entity.withinBounds(ps.xMinBound, ps.xMaxBound, ps.yMinBound, ps.yMaxBound) {
				ps.addEntity(entity)
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	test := PositionServer{xMinBound: 0, xMaxBound: 10, yMinBound: 0, yMaxBound: 10}
	test.createNewEntity()
	test.createNewEntity()
	fmt.Println(test.entities[0])
	test.tick()
	fmt.Println(test.entities[0])
	// test.entitiesInBounds(1, 2, 1, 2)
	// test.entitiesInBounds(11, 12, 11, 12)

	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	test.tick()
	fmt.Println(len(test.entities))
}
