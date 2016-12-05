package main

import (
	"fmt"
	"math/rand"
	"time"
)

var passedChanBufSize = 10

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println(6)
	fmt.Println(getXY(5, 4, 10))
	fmt.Println(7)
	fmt.Println(getXY(6, 4, 10))
	fmt.Println(11)
	fmt.Println(getXY(10, 4, 10))
	fmt.Println(16)
	fmt.Println(getXY(15, 4, 10))
	fmt.Println(1)
	fmt.Println(getXY(0, 4, 10))
	// fmt.Println(getXY(16, 4, 10))

	// test := PositionServer{xMinBound: 0, xMaxBound: 10, yMinBound: 0, yMaxBound: 10}
	// test.createNewEntity()
	// test.createNewEntity()
}
