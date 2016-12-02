package main

import "math/rand"

type PositionServer struct {
	xMinBound           int
	xMaxBound           int
	yMinBound           int
	yMaxBound           int
	color               string
	entities            []*Entity
	PassedEntityChannel chan (Entity)
	NewEntityChannel    chan (int)
	//AdjacentPS        AdjacentPositionServers
}

/*
type AdjacentPositionServers struct {
  left    chan(Entity)
  right   chan(Entity)
  top     chan(Entity)
  bottom  chan(Entity)
}
*/

func newPositionServer(xMinBound int, xMaxBound int, yMinBound int, yMaxBound int, color string) *PositionServer {
	return &PositionServer{
		xMinBound:           xMinBound,
		xMaxBound:           xMaxBound,
		yMinBound:           yMinBound,
		yMaxBound:           yMaxBound,
		color:               color,
		entities:            []*Entity{},
		PassedEntityChannel: make(chan Entity),
		NewEntityChannel:    make(chan int),
		//AdjacentPS:        AdjacentPositionServers{}
	}
}

// func (ps *PositionServer) addAdjacentPositionServers() {
//
// }

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
	addedEntity := &Entity{
		color: ps.color,
	}
	ps.entities = append(ps.entities, addedEntity)
}

func (ps *PositionServer) tick() {
	ps.moveEntities()
	ps.removeOutOfBoundsEntities()
	//lockstep
	ps.processPassedEntityChan()
	ps.processNewEntityChan()
	//send metrics
	//lockstep
}

func (ps *PositionServer) processNewEntityChan() {

}

func (ps *PositionServer) processPassedEntityChan() {

}

func (ps *PositionServer) removeOutOfBoundsEntities() {
	for i := len(ps.entities) - 1; i >= 0; i-- {
		entity := ps.entities[i]
		if !entity.withinBounds(ps.xMinBound, ps.xMaxBound, ps.yMinBound, ps.yMaxBound) {
			//TODO: move entity here
			if i != len(ps.entities)-1 {
				ps.entities = append(ps.entities[:i], ps.entities[i+1:]...)
			} else {
				ps.entities = ps.entities[:i]
			}
		}
	}
}

func (ps *PositionServer) moveEntities() {
	for _, entity := range ps.entities {
		entity.move()
	}
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
