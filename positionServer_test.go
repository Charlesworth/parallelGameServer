package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestPositionServer_setup(t *testing.T) {
	//As PositionServer uses random number gen, get a new seed for each test
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestPositionServer_tick(t *testing.T) {
	t.Error("remove this test")
	rand.Seed(time.Now().UTC().UnixNano())

	testPS := newPositionServer(0, 10, 0, 10, "")
	testPS.createNewEntity()
	testPS.createNewEntity()
	testPS.createNewEntity()

	for i := 11; i > 0; i-- {
		testPS.tick()
		for i, a := range testPS.entities {
			fmt.Println(i, *a)
		}
	}

	if len(testPS.entities) != 0 {
		t.Error()
	}
}

func TestPositionServer_removeOutOfBoundsEntities(t *testing.T) {
	t.Error("implement me")
}

func TestPositionServer_entitiesInBounds(t *testing.T) {
	testPS := newPositionServer(0, 10, 0, 10, "")

	testEntity := Entity{
		xPos:      2,
		yPos:      2,
		direction: 0,
	}

	testPS.entities = append(testPS.entities, &testEntity)

	returnedEntities := testPS.entitiesInBounds(0, 10, 0, 10)
	if len(returnedEntities) != 1 {
		t.Error()
	}

	returnedEntities = testPS.entitiesInBounds(5, 10, 5, 10)
	if len(returnedEntities) != 0 {
		t.Error()
	}

	returnedEntities = testPS.entitiesInBounds(-10, 10, -10, 10)
	if len(returnedEntities) != 1 {
		t.Error()
	}

	returnedEntities = testPS.entitiesInBounds(-10, -20, -10, -20)
	if len(returnedEntities) != 0 {
		t.Error()
	}
}

func TestPositionServer_intersectsBounds(t *testing.T) {
	testPS := newPositionServer(0, 10, 0, 10, "")

	if !testPS.intersectsBounds(0, 10, 0, 10) {
		t.Error()
	}

	if !testPS.intersectsBounds(-5, 5, -5, 5) {
		t.Error()
	}

	if !testPS.intersectsBounds(5, 15, 5, 15) {
		t.Error()
	}

	if testPS.intersectsBounds(-5, -15, 5, 15) {
		t.Error()
	}

	if testPS.intersectsBounds(5, 15, -5, -15) {
		t.Error()
	}

	if testPS.intersectsBounds(-1, -10, -1, -10) {
		t.Error()
	}

	if testPS.intersectsBounds(11, 20, 11, 20) {
		t.Error()
	}
}
