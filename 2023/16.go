package main

import "fmt"

var lines []string

// have we seen light in this tile passing in one of the 4 cardinal directions?
// 0 North
// 1 East
// 2 South
// 3 West
var energisedTiles [][][4]bool

func propagateLight(x, y int, direction int) {
	// if we are out of bounds, do nothing
	if x < 0 || y < 0 || y >= len(lines) || x >= len(lines[0]) {
		return
	}

	// have we seen this before? if so there is nothing to do
	if energisedTiles[y][x][direction] {
		return
	}

	// mark this as seen
	energisedTiles[y][x][direction] = true

	// process what we see
	switch lines[y][x] {
	case '.':
		switch direction {
		case 0:
			propagateLight(x, y-1, direction)
		case 1:
			propagateLight(x+1, y, direction)
		case 2:
			propagateLight(x, y+1, direction)
		case 3:
			propagateLight(x-1, y, direction)
		default:
			panic("invalid direction")
		}
	case '-':
		if direction == 0 || direction == 2 {
			propagateLight(x, y, 1)
			propagateLight(x, y, 3)
		} else if direction == 1 {
			propagateLight(x+1, y, direction)

		} else if direction == 3 {
			propagateLight(x-1, y, direction)

		} else {
			panic("invalid direction")
		}
	case '|':
		if direction == 1 || direction == 3 {
			propagateLight(x, y, 0)
			propagateLight(x, y, 2)
		} else if direction == 0 {
			propagateLight(x, y-1, direction)

		} else if direction == 2 {
			propagateLight(x, y+1, direction)

		} else {
			panic("invalid direction")
		}
	case '/':
		switch direction {
		case 0:
			propagateLight(x+1, y, 1)
		case 1:
			propagateLight(x, y-1, 0)
		case 2:
			propagateLight(x-1, y, 3)
		case 3:
			propagateLight(x, y+1, 2)
		default:
			panic("invalid direction, case /")
		}

	case '\\':
		switch direction {
		case 0:
			propagateLight(x-1, y, 3)
		case 1:
			propagateLight(x, y+1, 2)
		case 2:
			propagateLight(x+1, y, 1)
		case 3:
			propagateLight(x, y-1, 0)
		default:
			panic("invalid direction, case \\")
		}
	default:
		panic("invalid tile content")
	}

}

func resetState() {
	energisedTiles = make([][][4]bool, 0)
	for _ = range lines {
		energisedTiles = append(energisedTiles, make([][4]bool, len(lines[0])))
	}
}

func countSum() int {
	sum := 0
	for _, energisedRow := range energisedTiles {
		for _, energisedCell := range energisedRow {
			if energisedCell[0] || energisedCell[1] || energisedCell[2] || energisedCell[3] {
				sum++
			}
		}
	}
	return sum
}

func main() {
	lines = fetchLines(16)

	////////////////////////////////////////////////

	resetState()
	propagateLight(0, 0, 1)
	fmt.Println(countSum())

	////////////////////////////////////////////////

	maxSum := -1
	// try starting at top or bottom
	var sum int
	for x := range lines[0] {
		resetState()
		propagateLight(x, 0, 2)
		sum = countSum()
		if sum > maxSum {
			maxSum = sum
		}
		resetState()
		propagateLight(x, len(lines)-1, 0)
		sum = countSum()
		if sum > maxSum {
			maxSum = sum
		}
	}
	// try starting at left or right
	for y := range lines {
		resetState()
		propagateLight(0, y, 1)
		sum = countSum()
		if sum > maxSum {
			maxSum = sum
		}
		resetState()
		propagateLight(len(lines[0])-1, y, 3)
		sum = countSum()
		if sum > maxSum {
			maxSum = sum
		}
	}
	fmt.Println(maxSum)

}
