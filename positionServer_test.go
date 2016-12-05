package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestPositionServer_setup(t *testing.T) {
	//As PositionServer uses random number gen, get a new seed for each test
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestPositionServer_processNewEntityChannel(t *testing.T) {
	testPS := newPositionServer(0, 10, 0, 10, "")
	testPS.NewEntityChannel <- 2
	testPS.processNewEntityChan()

	if len(testPS.entities) != 2 {
		t.Error()
	}
}

func TestPositionServer_processPassedEntityChannel(t *testing.T) {
	testPS := newPositionServer(0, 10, 0, 10, "")
	testPS.PassedEntityChannel <- Entity{
		xPos:      5,
		yPos:      5,
		direction: 0,
	}
	testPS.processPassedEntityChan()

	if len(testPS.entities) != 1 {
		t.Error()
	}
}

func TestPositionServer_removeOutOfBoundsEntities(t *testing.T) {
	testPS := newPositionServer(0, 10, 0, 10, "")

	inBoundsEntity := Entity{
		xPos:      5,
		yPos:      5,
		direction: 0,
	}

	outOfBoundsEntity := Entity{
		xPos:      12,
		yPos:      12,
		direction: 0,
	}

	anotherOutOfBoundsEntity := Entity{
		xPos:      12,
		yPos:      12,
		direction: 0,
	}

	testPS.entities = append(testPS.entities, &outOfBoundsEntity)
	testPS.entities = append(testPS.entities, &inBoundsEntity)
	testPS.entities = append(testPS.entities, &anotherOutOfBoundsEntity)

	testPS.removeOutOfBoundsEntities()
	if len(testPS.entities) != 1 {
		t.Error()
	}
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

func TestPositionServer_waitForPassedEntities(t *testing.T) {
	//No need for error message here, will deadlock and error if fail
	testServer := newPositionServer(0, 0, 0, 0, "")
	//expect 4 confirmations on every sever
	testServer.PassedEntConfirmations <- true
	testServer.PassedEntConfirmations <- true
	testServer.PassedEntConfirmations <- true
	testServer.PassedEntConfirmations <- true
	if !testServer.waitForPassedEntities() {
		t.Error()
	}

	//expect 0 confirmations
	testServer.PassedEntConfirmations <- false
	testServer.PassedEntConfirmations <- false
	testServer.PassedEntConfirmations <- false
	testServer.PassedEntConfirmations <- false
	if testServer.waitForPassedEntities() {
		t.Error()
	}

	//expect 0 confirmations
	testServer.PassedEntConfirmations <- false
	testServer.PassedEntConfirmations <- false
	testServer.PassedEntConfirmations <- false
	testServer.PassedEntConfirmations <- true
	if !testServer.waitForPassedEntities() {
		t.Error()
	}
}
