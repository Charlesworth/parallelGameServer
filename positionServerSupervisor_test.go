package main

import "testing"

type getAdjacentTableTest struct {
	serverNumber    int
	totalServers    int
	squareRoot      int
	serverAdjacency serverAdjacency
}

var getAdjacentTests = []getAdjacentTableTest{
	getAdjacentTableTest{
		serverNumber: 0,
		totalServers: 16,
		squareRoot:   4,
		serverAdjacency: serverAdjacency{
			left:  3,
			right: 1,
			above: 12,
			below: 4,
		},
	},
	getAdjacentTableTest{
		serverNumber: 15,
		totalServers: 16,
		squareRoot:   4,
		serverAdjacency: serverAdjacency{
			left:  14,
			right: 12,
			above: 11,
			below: 3,
		},
	},
	getAdjacentTableTest{
		serverNumber: 4,
		totalServers: 9,
		squareRoot:   3,
		serverAdjacency: serverAdjacency{
			left:  3,
			right: 5,
			above: 1,
			below: 7,
		},
	},
	getAdjacentTableTest{
		serverNumber: 8,
		totalServers: 9,
		squareRoot:   3,
		serverAdjacency: serverAdjacency{
			left:  7,
			right: 6,
			above: 5,
			below: 2,
		},
	},
}

func Test_getAdjacentServerNumbers(t *testing.T) {
	for _, test := range getAdjacentTests {
		resultServerAd := getAdjacentServerNumbers(test.serverNumber, test.totalServers, test.squareRoot)
		if (resultServerAd.left != test.serverAdjacency.left) &&
			(resultServerAd.right != test.serverAdjacency.right) &&
			(resultServerAd.above != test.serverAdjacency.above) &&
			(resultServerAd.below != test.serverAdjacency.below) {
			t.Error("getAdjacentServerNumbers() error: serverNumber[", test.serverNumber, "] totalServers[", test.totalServers, "]")
		}
	}
}

type getXYCoordinatesTableTest struct {
	serverNumber int
	squareRoot   int
	length       int
	xMinResult   int
	xMaxResult   int
	yMinResult   int
	yMaxResult   int
}

var getXYCoordinatesTests = []getXYCoordinatesTableTest{
	getXYCoordinatesTableTest{
		serverNumber: 0,
		squareRoot:   4,
		length:       10,
		xMinResult:   0,
		xMaxResult:   10,
		yMinResult:   0,
		yMaxResult:   10,
	},
	getXYCoordinatesTableTest{
		serverNumber: 15,
		squareRoot:   4,
		length:       10,
		xMinResult:   30,
		xMaxResult:   40,
		yMinResult:   30,
		yMaxResult:   40,
	},
	getXYCoordinatesTableTest{
		serverNumber: 2,
		squareRoot:   3,
		length:       10,
		xMinResult:   20,
		xMaxResult:   30,
		yMinResult:   20,
		yMaxResult:   30,
	},
	getXYCoordinatesTableTest{
		serverNumber: 6,
		squareRoot:   3,
		length:       10,
		xMinResult:   0,
		xMaxResult:   10,
		yMinResult:   20,
		yMaxResult:   30,
	},
}

func Test_getXYCoordinates(t *testing.T) {
	for _, test := range getXYCoordinatesTests {
		xMinTest, xMaxTest, yMinTest, yMaxTest := getXYCoordinates(test.serverNumber, test.squareRoot, test.length)
		if (test.xMinResult != xMinTest) &&
			(test.xMaxResult != xMaxTest) &&
			(test.yMinResult != yMinTest) &&
			(test.yMaxResult != yMaxTest) {
			t.Error("getXYCoordinates() error: serverNumber[", test.serverNumber, "] totalServers[", test.squareRoot, "]")
		}
	}
}
