package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

var bufferSize = 100
var outOfBoundEntities chan (Entity)

func init() {
	outOfBoundEntities = make(chan Entity, bufferSize)
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

type PositionServer struct {
	xMinBound          int
	xMaxBound          int
	yMinBound          int
	yMaxBound          int
	color              string
	entities           []*Entity
	inputEntityChannel chan (Entity)
}

func (ps *PositionServer) createNewEntity() {
	xPos := ps.xMinBound + ((ps.xMaxBound - ps.xMinBound) / 2)
	yPos := ps.yMinBound + ((ps.yMaxBound - ps.yMinBound) / 2)
	newEntity := &Entity{
		xPos:      xPos,
		yPos:      yPos,
		direction: rand.Intn(4),
		color:     ps.color,
	}
	ps.entities = append(ps.entities, newEntity)
}

func (ps *PositionServer) addEntity(Entity) {
	newEntity := &Entity{
		color: ps.color,
	}
	ps.entities = append(ps.entities, newEntity)
}

func (ps *PositionServer) tick() {
	for i := len(ps.entities) - 1; i >= 0; i-- {
		entity := ps.entities[i]
		entity.move()
		if !entity.withinBounds(ps.xMinBound, ps.xMaxBound, ps.yMinBound, ps.yMaxBound) {
			log.Println("outOfBounds")
			outOfBoundEntities <- *entity
			if i != len(ps.entities)-1 {
				ps.entities = append(ps.entities[:i], ps.entities[i+1:]...)
			} else {
				ps.entities = ps.entities[:i]
			}
		}
	}
	// for _, entity := range ps.entities {
	//   entity.move()
	//   if !entity.withinBounds(ps.xMinBound, ps.xMaxBound, ps.yMinBound, ps.yMaxBound) {
	//     log.Println("outOfBounds")
	//     outOfBoundEntities <- *entity
	//
	//   }
	// }
}

func (ps *PositionServer) entitiesInBounds(xMinBound int, xMaxBound int, yMinBound int, yMaxBound int) []Entity {
	//First check that the input bounds overlaps the PositionServer's bounds, if not return empty entity slice
	if !ps.intersectsBounds(xMinBound, xMaxBound, yMinBound, yMaxBound) {
		return []Entity{}
	}

	//If it does intersect, range through all entities and return all within intersection
	entitiesInBounds := []Entity{}
	for _, entity := range ps.entities {
		if entity.withinBounds(xMinBound, xMaxBound, yMinBound, yMaxBound) {
			entitiesInBounds = append(entitiesInBounds, *entity)
		}
	}

	return entitiesInBounds
}

func (ps *PositionServer) intersectsBounds(xMinBound int, xMaxBound int, yMinBound int, yMaxBound int) bool {
	return !(ps.xMinBound > xMaxBound || xMinBound > ps.xMaxBound || ps.yMinBound > yMaxBound || yMinBound > ps.yMaxBound)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	test := PositionServer{xMinBound: 0, xMaxBound: 10, yMinBound: 0, yMaxBound: 10, inputEntityChannel: make(chan Entity, bufferSize)}
	test.createNewEntity()
	test.createNewEntity()
	fmt.Println(test.entities[0])
	test.tick()
	fmt.Println(test.entities[0])
	test.entitiesInBounds(1, 2, 1, 2)
	test.entitiesInBounds(11, 12, 11, 12)
}

type Entity struct {
	xPos      int
	yPos      int
	direction int
	color     string
}

func (e *Entity) move() {
	switch e.direction {
	case 0:
		e.xPos = e.xPos + 1
	case 1:
		e.xPos = e.xPos - 1
	case 2:
		e.yPos = e.yPos + 1
	case 3:
		e.yPos = e.yPos - 1
	default:
		log.Fatalln("direction input at actor.move() out of acceptable range")
	}
}

func (e *Entity) withinBounds(xMinBound int, xMaxBound int, yMinBound int, yMaxBound int) bool {
	withinXBounds := (e.xPos <= xMaxBound) && (e.xPos >= xMinBound)
	withinYBounds := (e.xPos <= xMaxBound) && (e.xPos >= xMinBound)
	return withinXBounds && withinYBounds
}
