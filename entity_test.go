package main

import "testing"

func TestEnitiy_move(t *testing.T) {
	testEntity := Entity{
		xPos:      0,
		yPos:      0,
		direction: 0,
	}

	testEntity.move()
	if testEntity.xPos != 1 {
		t.Error()
	}

	testEntity.direction = 1
	testEntity.move()
	if testEntity.xPos != 0 {
		t.Error()
	}

	testEntity.direction = 2
	testEntity.move()
	if testEntity.yPos != 1 {
		t.Error()
	}

	testEntity.direction = 3
	testEntity.move()
	if testEntity.yPos != 0 {
		t.Error()
	}
}

func TestEnitiy_withinBounds(t *testing.T) {
	testEntity := Entity{
		xPos:      5,
		yPos:      5,
		direction: 0,
	}

	if !testEntity.withinBounds(0, 10, 0, 10) {
		t.Error()
	}

	if testEntity.withinBounds(10, 20, 10, 20) {
		t.Error()
	}
}
