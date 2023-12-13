package main

import (
	"fmt"
	"strings"
)

var lines []string

// https://stackoverflow.com/a/24893667/14464173
func replaceAtIndex(input string, replacement byte, index int) string {
	return strings.Join([]string{input[:index], string(replacement), input[index+1:]}, "")
}

func moveDownPipe(yCurrent, xCurrent, yLast, xLast int) (int, int) {
	switch lines[yCurrent][xCurrent] {
	case '|':
		if yLast == yCurrent-1 {
			return yCurrent + 1, xCurrent
		} else {
			return yCurrent - 1, xCurrent
		}
	case '-':
		if xLast == xCurrent-1 {
			return yCurrent, xCurrent + 1
		} else {
			return yCurrent, xCurrent - 1
		}
	case 'L':
		if yLast == yCurrent-1 {
			return yCurrent, xCurrent + 1
		} else {
			return yCurrent - 1, xCurrent
		}
	case 'J':
		if yLast == yCurrent-1 {
			return yCurrent, xCurrent - 1
		} else {
			return yCurrent - 1, xCurrent
		}
	case '7':
		if yLast == yCurrent+1 {
			return yCurrent, xCurrent - 1
		} else {
			return yCurrent + 1, xCurrent
		}
	case 'F':
		if yLast == yCurrent+1 {
			return yCurrent, xCurrent + 1
		} else {
			return yCurrent + 1, xCurrent
		}

	default:
		panic("unexpected tile rune")
	}
}

func main() {
	lines = fetchLines(10)

	////////////////////////////////////////////////

	// these contain will contain two coordinates each as we move through the tunnels
	forwardPointer := make([]int, 2)
	backwardPointer := make([]int, 2)
	var yStart, xStart int

outer:
	for y := range lines {
		for x := range lines[y] {
			if lines[y][x] == 'S' {
				forwardPointer[0] = y
				backwardPointer[0] = y
				forwardPointer[1] = x
				backwardPointer[1] = x
				yStart = y
				xStart = x
				break outer
			}
		}
	}

	sIndicator := 0

	// assume we don't have to do bounds checking and that S isn't on the very edge of the grid
	if lines[yStart-1][xStart] == '|' || lines[yStart-1][xStart] == '7' || lines[yStart-1][xStart] == 'F' {
		// up
		forwardPointer = append(forwardPointer, yStart-1)
		forwardPointer = append(forwardPointer, xStart)
		sIndicator += 1000
	}
	if lines[yStart+1][xStart] == '|' || lines[yStart+1][xStart] == 'L' || lines[yStart+1][xStart] == 'F' {
		// down
		if len(forwardPointer) == 2 {
			forwardPointer = append(forwardPointer, yStart+1)
			forwardPointer = append(forwardPointer, xStart)
		} else {
			backwardPointer = append(backwardPointer, yStart+1)
			backwardPointer = append(backwardPointer, xStart)
		}
		sIndicator += 100
	}
	if lines[yStart][xStart-1] == '-' || lines[yStart][xStart-1] == 'L' || lines[yStart][xStart-1] == 'F' {
		// left
		if len(forwardPointer) == 2 {
			forwardPointer = append(forwardPointer, yStart)
			forwardPointer = append(forwardPointer, xStart-1)
		} else {
			backwardPointer = append(backwardPointer, yStart)
			backwardPointer = append(backwardPointer, xStart-1)
		}
		sIndicator += 10
	}
	if lines[yStart][xStart+1] == '-' || lines[yStart][xStart+1] == '7' || lines[yStart][xStart+1] == 'J' {
		// right
		backwardPointer = append(backwardPointer, yStart)
		backwardPointer = append(backwardPointer, xStart+1)
		sIndicator += 1
	}

	switch sIndicator {
	case 1100:
		lines[yStart] = replaceAtIndex(lines[yStart], '|', xStart)
	case 11:
		lines[yStart] = replaceAtIndex(lines[yStart], '-', xStart)
	case 1001:
		lines[yStart] = replaceAtIndex(lines[yStart], 'L', xStart)
	case 1010:
		lines[yStart] = replaceAtIndex(lines[yStart], 'J', xStart)
	case 110:
		lines[yStart] = replaceAtIndex(lines[yStart], '7', xStart)
	case 101:
		lines[yStart] = replaceAtIndex(lines[yStart], 'F', xStart)
	default:
		panic("unexpected S state")
	}

	if !(len(forwardPointer) == 4 && len(backwardPointer) == 4) {
		panic("unexpected initial pointer state")
	}

	// (data structures for part 2)
	// all the tiles in the loop
	inLoop := make(map[[2]int]bool)
	inLoop[[2]int{forwardPointer[0], forwardPointer[1]}] = true
	inLoop[[2]int{forwardPointer[2], forwardPointer[3]}] = true
	yCurrent, xCurrent := forwardPointer[2], forwardPointer[3]
	var yNext, xNext int

	var y, x int
	for count := 1; true; count++ {
		// the heads of both pointers hit the same tile
		if forwardPointer[2] == backwardPointer[2] && forwardPointer[3] == backwardPointer[3] {
			fmt.Println(count)
			break
		}

		// the pointers have moved past each other
		if forwardPointer[2] == backwardPointer[0] && forwardPointer[3] == backwardPointer[1] && backwardPointer[2] == forwardPointer[0] && backwardPointer[3] == forwardPointer[1] {
			fmt.Println(count - 1)
			break
		}

		// move the pointers along
		y, x = moveDownPipe(forwardPointer[2], forwardPointer[3], forwardPointer[0], forwardPointer[1])
		forwardPointer = append(forwardPointer, y)
		forwardPointer = append(forwardPointer, x)
		forwardPointer = forwardPointer[2:]
		y, x = moveDownPipe(backwardPointer[2], backwardPointer[3], backwardPointer[0], backwardPointer[1])
		backwardPointer = append(backwardPointer, y)
		backwardPointer = append(backwardPointer, x)
		backwardPointer = backwardPointer[2:]
	}

	////////////////////////////////////////////////

	// populate the inLoop map
	for {
		// move along
		yNext, xNext = moveDownPipe(yCurrent, xCurrent, yStart, xStart)

		// if we have reached the start, break out
		if inLoop[[2]int{yNext, xNext}] {
			break
		}

		// otherwise set this next cell as visited
		inLoop[[2]int{yNext, xNext}] = true

		yCurrent, xCurrent, yStart, xStart = yNext, xNext, yCurrent, xCurrent
	}

	innerTiles := 0
	for y := range lines {
		leftIsInner := false
		rightIsInner := false
		for x := range lines[y] {
			tile := [2]int{y, x}

			// if this tile is not part of the main loop
			if !inLoop[tile] {
				// if both left and right are inner, this whole tile is inner
				if leftIsInner && rightIsInner {
					innerTiles += 1
				}
				// in any case, there is nothing more to do
				continue
			}

			// now we know that this tile makes up the main loop
			switch lines[y][x] {
			case '|':
				leftIsInner = !leftIsInner
				rightIsInner = !rightIsInner
			case '-':
			case 'L':
				leftIsInner = !leftIsInner
			case 'J':
				leftIsInner = !leftIsInner
			case '7':
				rightIsInner = !rightIsInner
			case 'F':
				rightIsInner = !rightIsInner
			default:
				panic("unexpected tile rune")
			}

		}
	}

	fmt.Println(innerTiles)
}
