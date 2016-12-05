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
