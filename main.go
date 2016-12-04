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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	test := &PositionServerSupervisor{}
	for i := 0; i < 26; i++ {
		err := test.initServers(i, 0, 0)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(i)
		}
	}

	// test := PositionServer{xMinBound: 0, xMaxBound: 10, yMinBound: 0, yMaxBound: 10}
	// test.createNewEntity()
	// test.createNewEntity()
}
