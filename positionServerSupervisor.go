package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

type PositionServerSupervisor struct {
	positionServers []*PositionServer
	waitGroup       sync.WaitGroup
}

func newPositionServerSupervisor() *PositionServerSupervisor {
	return &PositionServerSupervisor{
		positionServers: []*PositionServer{},
		waitGroup:       sync.WaitGroup{},
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

	//for numberOfPositionServers
	for i := 0; i < numberOfPositionServers; i++ {
		//find bounds and init server
		xMin, xMax, yMin, yMax := getXYCoordinates(i, squareRoot, sideLengthPerServer)
		newPS := newPositionServer(xMin, xMax, yMin, yMax, "red", pss.waitGroup) //TODO color

		//add new entities
		for iNewEntities := startingEntitiesPerServer; iNewEntities > 0; iNewEntities-- {
			newPS.createNewEntity()
		}
	}

	//pass each server its adjacent servers' PassEntity channels
	for serverNumber := range pss.positionServers {
		adjacentServersNumbers := getAdjacentServerNumbers(serverNumber, numberOfPositionServers, squareRoot)
		pss.referenceAdjacentPassedEntityChannels(serverNumber, adjacentServersNumbers)
	}

	return nil
}

func (pss *PositionServerSupervisor) startServers() {
	pss.waitGroup.Add(len(pss.positionServers))

	for _, positionServer := range pss.positionServers {
		go positionServer.mainLoop()
	}

	//TODO potential timing issue when wait is freed, if this thread goes last then
	//no wait will be set
	for {
		pss.waitGroup.Wait()
		pss.waitGroup.Add(len(pss.positionServers))
	}

	/*
	  - send any newEntity() calls
	*/
}

type serverAdjacency struct {
	left  int
	right int
	above int
	below int
}

func (pss *PositionServerSupervisor) referenceAdjacentPassedEntityChannels(serverNo int, adjacentServerNumbers serverAdjacency) {
	pss.positionServers[serverNo].AdjacentPassChannels = AdjacentPassedEntChannels{
		leftPEChan:   pss.positionServers[adjacentServerNumbers.left].PassedEntityChannel,
		leftConfirm:  pss.positionServers[adjacentServerNumbers.left].PassedEntConfirmations,
		rightPEChan:  pss.positionServers[adjacentServerNumbers.right].PassedEntityChannel,
		rightConfirm: pss.positionServers[adjacentServerNumbers.right].PassedEntConfirmations,
		abovePEChan:  pss.positionServers[adjacentServerNumbers.above].PassedEntityChannel,
		aboveConfirm: pss.positionServers[adjacentServerNumbers.above].PassedEntConfirmations,
		belowPEChan:  pss.positionServers[adjacentServerNumbers.below].PassedEntityChannel,
		belowConfirm: pss.positionServers[adjacentServerNumbers.below].PassedEntConfirmations,
	}
}

//-------Helper Functions---------------

func getAdjacentServerNumbers(serverNumber int, totalServers int, squareRoot int) serverAdjacency {
	position := serverNumber + 1
	thisServer := serverAdjacency{}

	//LEFT
	if position%squareRoot == 1 {
		thisServer.left = position + (squareRoot - 1)
	} else {
		thisServer.left = position - 1
	}

	//RIGHT
	if position%squareRoot == 0 {
		thisServer.right = position - (squareRoot - 1)
	} else {
		thisServer.right = position + 1
	}

	//ABOVE
	if position <= squareRoot {
		thisServer.above = totalServers + (-squareRoot + position)
	} else {
		thisServer.above = position - squareRoot
	}

	//BELOW
	if position > (totalServers - squareRoot) {
		thisServer.below = position - (totalServers - squareRoot)
	} else {
		thisServer.below = position + squareRoot
	}

	//Minus 1 from each field to as the position starts at 0
	thisServer = serverAdjacency{
		thisServer.left - 1,
		thisServer.right - 1,
		thisServer.above - 1,
		thisServer.below - 1}

	return thisServer
}

func getXYCoordinates(i int, squareRoot int, length int) (xMin int, xMax int, yMin int, yMax int) {
	i = i + 1

	var row int
	fltI, fltSquareRoot := float64(i), float64(squareRoot)
	if (fltI / fltSquareRoot) == math.Trunc(fltI/fltSquareRoot) {
		row = (i / squareRoot) - 1
	} else {
		row = i / squareRoot
	}

	var col int
	if i%squareRoot == 0 {
		col = squareRoot - 1
	} else {
		col = (i % squareRoot) - 1
	}

	xMin = col * length
	xMax = (col + 1) * length
	yMin = row * length
	yMax = (row + 1) * length
	return xMin, xMax, yMin, yMax
}
