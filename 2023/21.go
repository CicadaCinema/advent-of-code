package main

import "fmt"

func takeStep(startCoords map[[2]int]bool, allowedCells [][]bool) map[[2]int]bool {
	result := make(map[[2]int]bool)

	for start, v := range startCoords {
		if !v {
			continue
		}

		// try N
		if start[0] > 0 && allowedCells[start[0]-1][start[1]] {
			result[[2]int{start[0] - 1, start[1]}] = true
		}
		// try S
		if start[0] < len(allowedCells)-1 && allowedCells[start[0]+1][start[1]] {
			result[[2]int{start[0] + 1, start[1]}] = true
		}
		// try W
		if start[1] > 0 && allowedCells[start[0]][start[1]-1] {
			result[[2]int{start[0], start[1] - 1}] = true
		}
		// try E
		if start[1] < len(allowedCells[0])-1 && allowedCells[start[0]][start[1]+1] {
			result[[2]int{start[0], start[1] + 1}] = true
		}
	}

	return result
}

func countPositionsFiniteGrid(startX, startY, stepCount int, allowedCells [][]bool) int {
	allowedPositions := make(map[[2]int]bool)
	allowedPositions[[2]int{startY, startX}] = true

	for step := 1; step <= stepCount; step++ {
		allowedPositions = takeStep(allowedPositions, allowedCells)
		//fmt.Printf("%d steps can get you to %d places\n", step, len(allowedPositions))
	}

	return len(allowedPositions)
}

// starting at a particular point on a finite grid, detect a repetitive number of cells reached
// this function returns an output a, b, c such that:
// stepping >=a steps, where the number of steps is even, can make us reach b spots
// stepping >=a steps, where the number of steps is odd, can make us reach c spots
func detectCycle(startX, startY int, allowedCells [][]bool) (int, int, int) {
	// step counter
	stepsTaken := 0

	allowedPositions := make(map[[2]int]bool)
	allowedPositions[[2]int{startY, startX}] = true

	placesReached := make([]int, 0)

	// warm up, there will definitely be no cycle here
	for step := 0; step < 100; step++ {
		allowedPositions = takeStep(allowedPositions, allowedCells)
		stepsTaken++
		placesReached = append(placesReached, len(allowedPositions))
	}

	var lastIndex int
	var evenCount, oddCount int

	for {
		allowedPositions = takeStep(allowedPositions, allowedCells)
		stepsTaken++
		placesReached = append(placesReached, len(allowedPositions))

		// detect threefold repetition
		lastIndex = len(placesReached) - 1

		a := placesReached[lastIndex] == placesReached[lastIndex-2]
		b := placesReached[lastIndex-2] == placesReached[lastIndex-4]
		c := placesReached[lastIndex-1] == placesReached[lastIndex-3]
		d := placesReached[lastIndex-3] == placesReached[lastIndex-5]
		if a && b && c && d {
			if stepsTaken%2 == 0 {
				// we have taken an even number of steps
				evenCount = placesReached[lastIndex]
				oddCount = placesReached[lastIndex-1]
			} else {
				// we have taken an odd number of steps
				evenCount = placesReached[lastIndex-1]
				oddCount = placesReached[lastIndex]
			}
			break
		}
	}

	// assert for 200 more steps
	for step := 0; step < 200; step++ {
		allowedPositions = takeStep(allowedPositions, allowedCells)
		stepsTaken++

		if stepsTaken%2 == 0 {
			assert(len(allowedPositions) == evenCount, "evenCount assertion broken")
		} else {
			assert(len(allowedPositions) == oddCount, "evenCount assertion broken")
		}
	}

	return lastIndex, evenCount, oddCount
}

func repeatHorizonal(grid [][]bool, n int) [][]bool {
	assert(n >= 1, "invalid input")

	result := make([][]bool, 0)

	for _, row := range grid {
		// duplicate this row n times
		newRow := make([]bool, 0)
		for i := 0; i < n; i++ {
			for _, cell := range row {
				newRow = append(newRow, cell)
			}
		}
		assert(len(newRow) == n*len(row), "invalid new row length")
		result = append(result, newRow)
	}

	assert(len(result) == len(grid), "invalid number of rows")
	return result
}

func repeatVertical(grid [][]bool, n int) [][]bool {
	assert(n >= 1, "invalid input")

	result := make([][]bool, 0)

	for i := 0; i < n; i++ {
		for _, row := range grid {
			newRow := make([]bool, len(row))
			copy(newRow, row)
			assert(len(newRow) == len(grid[0]), "invalid row length")
			result = append(result, newRow)
		}
	}

	assert(len(result) == n*len(grid), "invalid number of rows")

	return result
}

var dim, lowestIndex, highestIndex, middleIndex, reachableCells, stepsToTake int

func main() {
	lines := fetchLines(21)

	////////////////////////////////////////////////

	var startX, startY int
	foundStart := false

	// can walk here
	allowedCellsOriginal := make([][]bool, len(lines))
	for i := range allowedCellsOriginal {
		row := make([]bool, len(lines[0]))
		allowedCellsOriginal[i] = row
	}
	// at this point, nothing is allowed

	for i := range lines {
		for j := range lines[i] {
			if lines[i][j] == '.' {
				// this spot is allowed
				allowedCellsOriginal[i][j] = true
			} else if lines[i][j] == 'S' {
				allowedCellsOriginal[i][j] = true
				assert(!foundStart, "found two starts")
				startY = i
				startX = j
			}
		}
	}

	fmt.Println(countPositionsFiniteGrid(startX, startY, 64, allowedCellsOriginal))
	fmt.Printf("startX: %d\nstartY: %d\n", startX, startY)

	////////////////////////////////////////////////

	// the dimensions are 131 x 131
	// the lowest index is 0,
	// the highest index is 130
	// the middle index is 65
	dim = 131
	lowestIndex = 0
	highestIndex = 130
	middleIndex = 65

	reachableCells = 0

	stepsToTake = 26501365

	fmt.Println(reachableCells)
}
