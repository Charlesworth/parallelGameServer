package main

type Entity struct {
	xPos      int
	yPos      int
	direction int
	color     string
}

func (e *Entity) move() {
	switch e.direction {
	case 0:
		e.xPos = e.xPos + 1
	case 1:
		e.xPos = e.xPos - 1
	case 2:
		e.yPos = e.yPos + 1
	case 3:
		e.yPos = e.yPos - 1
	}
}

func (e *Entity) withinBounds(xMinBound int, xMaxBound int, yMinBound int, yMaxBound int) bool {
	withinXBounds := (e.xPos <= xMaxBound) && (e.xPos >= xMinBound)
	withinYBounds := (e.yPos <= yMaxBound) && (e.yPos >= yMinBound)
	return withinXBounds && withinYBounds
}
