package main

import (
	"log"
	"math/rand"
	"sync"
)

type PositionServer struct {
	xMinBound              int
	xMaxBound              int
	yMinBound              int
	yMaxBound              int
	color                  string
	entities               []*Entity
	NewEntityChannel       chan (int)
	PassedEntityChannel    chan (Entity)
	PassedEntConfirmations chan (bool)
	AdjacentPassChannels   AdjacentPassedEntChannels
	lockStepWG             *sync.WaitGroup
}

type AdjacentPassedEntChannels struct {
	leftPEChan   chan (Entity)
	leftConfirm  chan (bool)
	rightPEChan  chan (Entity)
	rightConfirm chan (bool)
	abovePEChan  chan (Entity)
	aboveConfirm chan (bool)
	belowPEChan  chan (Entity)
	belowConfirm chan (bool)
}

func newPositionServer(xMinBound int, xMaxBound int, yMinBound int, yMaxBound int, color string, lockStepWG *sync.WaitGroup) *PositionServer {
	return &PositionServer{
		xMinBound: xMinBound,
		xMaxBound: xMaxBound,
		yMinBound: yMinBound,
		yMaxBound: yMaxBound,
		color:     color,
		entities:  []*Entity{},
		//buffered, length of 1 (non blocking channel)
		NewEntityChannel: make(chan int, 1),
		//buffered, length of passedChanBufSize
		PassedEntityChannel: make(chan Entity, passedChanBufSize),
		//buffered, length of 4
		PassedEntConfirmations: make(chan bool, 4),
		lockStepWG:             lockStepWG,
	}
}

func (ps *PositionServer) mainLoop() {
	// for {
	ps.moveEntities()

	//check for outOfBounds Entities and send any to their new servers
	outOfBoundsEnitites := ps.removeOutOfBoundsEntities()
	if len(outOfBoundsEnitites) != 0 {
		ps.sendOutOfBoundEntities(outOfBoundsEnitites)
	} else {
		ps.confirmNonePassed()
	}

	//wait for their confirmation of passedEntitys send
	passedEnitities := ps.waitForPassedEntities()

	//if any entities where passed, process them
	if passedEnitities {
		ps.processPassedEntityChan()
	}

	//ps.processNewEntityChan() ERROR HERE
	//send metrics
	//render
	// if verbose {
	// ps.verboseLogs()
	// }
	ps.lockStep()
	// }
}

func (ps *PositionServer) verboseLogs() {
	log.Println("[x:", ps.xMinBound, ", y:", ps.yMinBound, "]",
		"\nentity number: ", len(ps.entities))
	for _, entity := range ps.entities {
		log.Println(entity)
	}
	log.Println("----------------------------------------------")
}

func (ps *PositionServer) confirmNonePassed() {
	ps.AdjacentPassChannels.leftConfirm <- false
	ps.AdjacentPassChannels.rightConfirm <- false
	ps.AdjacentPassChannels.aboveConfirm <- false
	ps.AdjacentPassChannels.belowConfirm <- false
}

func (ps *PositionServer) lockStep() {
	ps.lockStepWG.Done()
}

func (ps *PositionServer) waitForPassedEntities() (entitiesToProcess bool) {
	entitiesToProcess = false
	for sentConfirmations := 0; sentConfirmations < 4; {
		containsPassedEntity := <-ps.PassedEntConfirmations
		if containsPassedEntity {
			entitiesToProcess = true
		}
		sentConfirmations = sentConfirmations + 1
	}
	return
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

func (ps *PositionServer) addEntity(e Entity) {
	addedEntity := &Entity{
		xPos:      e.xPos,
		yPos:      e.yPos,
		direction: e.direction,
		color:     ps.color,
	}
	ps.entities = append(ps.entities, addedEntity)
}

func (ps *PositionServer) sendOutOfBoundEntities(entities []Entity) {
	leftConfirm, rightConfirm, aboveConfirm, belowConfirm := false, false, false, false

	for _, entity := range entities {
		if entity.xPos < ps.xMinBound {
			ps.AdjacentPassChannels.leftPEChan <- entity
			leftConfirm = true
			continue
		}
		if entity.xPos > ps.xMaxBound {
			ps.AdjacentPassChannels.rightPEChan <- entity
			rightConfirm = true
			continue
		}
		if entity.yPos < ps.yMinBound {
			ps.AdjacentPassChannels.abovePEChan <- entity
			aboveConfirm = true
			continue
		}
		if entity.yPos > ps.yMaxBound {
			ps.AdjacentPassChannels.belowPEChan <- entity
			belowConfirm = true
			continue
		}
	}

	ps.AdjacentPassChannels.leftConfirm <- leftConfirm
	ps.AdjacentPassChannels.rightConfirm <- rightConfirm
	ps.AdjacentPassChannels.aboveConfirm <- aboveConfirm
	ps.AdjacentPassChannels.belowConfirm <- belowConfirm
}

func (ps *PositionServer) processNewEntityChan() {
	iNewEntities := <-ps.NewEntityChannel
	for i := iNewEntities; i > 0; i-- {
		ps.createNewEntity()
	}
}

func (ps *PositionServer) processPassedEntityChan() {
	for {
		select {
		case entity := <-ps.PassedEntityChannel:
			log.Println("recieved:", entity)
			if !entity.withinBounds(ps.xMinBound, ps.xMaxBound, ps.yMinBound, ps.yMaxBound) {
				entity = ps.convertCircularPassedEntity(entity)
				log.Println("converted:", entity)
			}
			ps.addEntity(entity)
		default:
			return
		}
	}
}

func (ps *PositionServer) convertCircularPassedEntity(entity Entity) Entity {
	if entity.xPos > ps.xMaxBound {
		entity.xPos = ps.xMinBound
		return entity
	}
	if entity.xPos < ps.xMinBound {
		entity.xPos = ps.xMaxBound
		return entity
	}
	if entity.yPos > ps.yMaxBound {
		entity.yPos = ps.yMinBound
		return entity
	}
	if entity.yPos < ps.yMinBound {
		entity.yPos = ps.yMaxBound
		return entity
	}
	return entity
}

func (ps *PositionServer) removeOutOfBoundsEntities() (removedEntities []Entity) {
	for i := len(ps.entities) - 1; i >= 0; i-- {
		entity := ps.entities[i]
		if !entity.withinBounds(ps.xMinBound, ps.xMaxBound, ps.yMinBound, ps.yMaxBound) {
			removedEntities = append(removedEntities, *entity)
			if i != len(ps.entities)-1 {
				ps.entities = append(ps.entities[:i], ps.entities[i+1:]...)
			} else {
				ps.entities = ps.entities[:i]
			}
		}
	}
	return
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
