package main

import (
	"errors"
	"fmt"
	"math"
)

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

func (pss *PositionServerSupervisor) initServers(numberOfPositionServers int, sideLengthPerServer int, startingEntitiesPerServer int) error {
	//check number of servers is at least 4
	if numberOfPositionServers < 4 {
		return errors.New("PositionServerSupervisor.initServers() numberOfPositionServers requires a minimum of 4")
	}

	//check the the square root of numberOfPositionServers is a whole number
	squareRootFloat := math.Sqrt(float64(25))
	if math.Trunc(squareRootFloat) != squareRootFloat {
		return errors.New("PositionServerSupervisor.initServers() numberOfPositionServers must be a square of a whole number")
	}

	//check length of sides it at least 10
	if sideLengthPerServer < 10 {
		return errors.New("PositionServerSupervisor.initServers() sideLengthPerServer must be 10 or more")
	}

	//check that starting entities is at least 1
	if startingEntitiesPerServer < 1 {
		return errors.New("PositionServerSupervisor.initServers() startingEntitiesPerServer must be 1 or more")
	}

	squareRoot := int(squareRootFloat)
	fmt.Println(squareRoot)
	/*
	   for numberOfPositionServers
	     make each positionServer
	     for startingEntitiesPerServer
	       positionServer.newEntity
	   pass each server its adjacent servers' PassedEntityChannel
	*/

	return nil
}

func (pss *PositionServerSupervisor) referenceAdjacentPassedEntityChannels(sa serverAdjacency) {
	pss.positionServers[sa.serverNumber].AdjacentPassChannels = AdjacentPassedEntChannels{
		leftPEChan:   pss.positionServers[sa.left].PassedEntityChannel,
		leftConfirm:  pss.positionServers[sa.left].PassedEntConfirmations,
		rightPEChan:  pss.positionServers[sa.right].PassedEntityChannel,
		rightConfirm: pss.positionServers[sa.right].PassedEntConfirmations,
		abovePEChan:  pss.positionServers[sa.above].PassedEntityChannel,
		aboveConfirm: pss.positionServers[sa.above].PassedEntConfirmations,
		belowPEChan:  pss.positionServers[sa.below].PassedEntityChannel,
		belowConfirm: pss.positionServers[sa.below].PassedEntConfirmations,
	}
}

type serverAdjacency struct {
	serverNumber int
	left         int
	right        int
	above        int
	below        int
}

func getAdjacent(serverNumber int, totalServers int, squareRoot int) serverAdjacency {
	thisServer := serverAdjacency{serverNumber, 0, 0, 0, 0}

	//LEFT
	if serverNumber%squareRoot == 1 {
		thisServer.left = serverNumber + (squareRoot - 1)
	} else {
		thisServer.left = serverNumber - 1
	}

	//RIGHT
	if serverNumber%squareRoot == 0 {
		thisServer.right = serverNumber - (squareRoot - 1)
	} else {
		thisServer.right = serverNumber + 1
	}

	//above
	if serverNumber <= squareRoot {
		thisServer.above = totalServers + (-squareRoot + serverNumber)
	} else {
		thisServer.above = serverNumber - squareRoot
	}

	//BOTTOM
	if serverNumber > (totalServers - squareRoot) {
		thisServer.below = serverNumber - (totalServers - squareRoot)
	} else {
		thisServer.below = serverNumber + squareRoot
	}

	return thisServer
}

func (pss *PositionServerSupervisor) startServers() {
	for _, positionServer := range pss.positionServers {
		go positionServer.mainLoop()
	}

	/*
	  - supervise the wait group
	  - send any newEntity() calls
	*/
}
