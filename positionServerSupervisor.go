package main

import (
	"errors"
	"fmt"
	"math"
)

type PositionServerSupervisor struct {
	positionServers []*PositionServer
	//sync.waitgroup
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

	for i := 0; i < numberOfPositionServers; i++ {
		/*
		   for numberOfPositionServers
		     make each positionServer
		     for startingEntitiesPerServer
		       positionServer.newEntity
		*/
		//find next bounds
		// newPS := newPositionServer(i, xMinBound int, xMaxBound int, yMinBound int, yMaxBound int, "red")

		newPS := newPositionServer(0, 10, 0, 10, "red") //problem here, get proper bound input
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

func getXY(i int, squareRoot int, length int) (xMin int, xMax int, yMin int, yMax int) {
	i = i + 1

	/*
		  Fix this bit
			var row int
			if i / squareRoot == squareRoot
			row := i / squareRoot
			fmt.Println("row", row)
	*/
	row := i / squareRoot

	var col int
	if i%squareRoot == 0 {
		col = squareRoot - 1
	} else {
		col = (i % squareRoot) - 1
	}
	fmt.Println("col", col)

	xMin = col * length
	xMax = (col + 1) * length
	yMin = row * length
	yMax = (row + 1) * length
	return xMin, xMax, yMin, yMax
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

type serverAdjacency struct {
	left  int
	right int
	above int
	below int
}

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

	//above
	if position <= squareRoot {
		thisServer.above = totalServers + (-squareRoot + position)
	} else {
		thisServer.above = position - squareRoot
	}

	//BOTTOM
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
