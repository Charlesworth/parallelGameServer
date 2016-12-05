package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var passedChanBufSize = 10

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var servers int
	flag.IntVar(&servers, "servers", 3, "amount of position servers, must be a perfect square")

	var startingEntities int
	flag.IntVar(&startingEntities, "startingEntities", 1, "amount of entities per position server on start, minimum of 1")

	var sideLength int
	flag.IntVar(&sideLength, "sideLength", 10, "length of side of each server")

	flag.Parse()

	fmt.Println("Distributed Position Servers")
	fmt.Println("Starting with servers[", servers, "], startingEntities[", startingEntities, "], sideLength[", sideLength, "]")

	psSupervisor := newPositionServerSupervisor()
	err := psSupervisor.initServers(servers, sideLength, startingEntities)
	if err != nil {
		log.Fatal(err)
	}

	psSupervisor.startServers()

}
