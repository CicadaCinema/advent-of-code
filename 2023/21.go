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

/*

// UNUSED FUNCTIONS FOR PREVIOUS ATTEMPTS TO SOLVE THIS PROBLEM

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

*/

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

	////////////////////////////////////////////////

	// do some pre-processing to find areas completely enclosed by hashes (walls/rocks)
	// this is a reasonably fast heuristic because the grid is usually quite sparse,
	// so a wall of rocks enclosing a large area should be rare
	for y := range allowedCellsOriginal {
		for x := range allowedCellsOriginal[0] {
			if !allowedCellsOriginal[y][x] {
				continue
			}
			count := countPositionsFiniteGrid(x, y, 10, allowedCellsOriginal)
			if count < 10 {
				// for this "bad" enclosed cell/garden plot, just turn it into a rock
				allowedCellsOriginal[y][x] = false
			}
		}
	}

	////////////////////////////////////////////////

	/*
	   In this solution, rocks are called hashes.

	   Imagine that the cells of the infinite grid are coloured black and white, like a chessboard.
	   Denote by "parity" the colour of a cell.
	   For example, going between two cells of different parities always requires an odd number of steps.

	   The number of allowed steps is odd, so the only reachable cells will be the ones with a different parity to the starting cell.

	   In the absense of any hashes, the reachable region will be diamond-shaped, and the above condition means that roughly half of the cells in this diamond will be reachable.
	   This value of reachable cells is calculated by `diamondUnobstructedReachable`.

	   Note that we eliminated any awkward enclosed cells/garden plots in the preprocessing step above.

	   Also note that the number of steps is exactly of the form `n*131+65` (with `n=202300` in this case), and the presense of the diagonal free cells/garden plots in the input, as well as the nice horizontal/vertical rows of garden plots in the input, mean that we don't have to worry about what happens on the edge of the above-mentioned diamond. The edge of this diamond is away from the hashes.

	   The rest of the solution focuses on counting the number of hashes within the diamond, which we have mistakenly over-counted above as 'reachable cells'.

	   For this I drew a picture and worked on paper most of the time.

	   Note that when moving 131 steps from one grid-square to the next (there are infinite grids stretching out in all directions), we spend exactly 131 steps, an odd number. So the parity of the number of *remaining* steps changes.
	*/

	// the dimensions are 131 x 131
	// the lowest index is 0,
	// the highest index is 130
	// the middle index is 65
	dim := 131
	middleIndex := 65

	// allowed steps to take
	steps := 26501365

	diamondUnobstructedReachable := steps*(steps+1) + steps + 1

	// if one grid-step is 131 normal steps (or alternatively, picture moving from the centre of one grid to the centre of an adjacent grid), then this is the number of grid-steps we can take
	gridSteps := steps / dim

	// the number of hashes in the grid
	hashes := 0
	// the number of hashes in the grid, which are reachable by an even number of steps
	hashesEven := 0
	// the number of hashes in the diamond shape (in the input, not the large diamond described in the block comment above)
	diamondH := 0
	// the number of hashes in the diamond shape, which are reachable by an even number of steps
	diamondHOdd := 0
	for y := range allowedCellsOriginal {
		for x := range allowedCellsOriginal[0] {
			if !allowedCellsOriginal[y][x] {
				hashes++
				// the initial parity is 0 since x==y at the startig position
				if (y+x)%2 == 0 {
					hashesEven++
				}

				deltaY, deltaX := y-middleIndex, x-middleIndex
				deltaY = max(deltaY, -deltaY)
				deltaX = max(deltaX, -deltaX)
				// if we are inside the diamond
				if deltaX+deltaY < middleIndex-1 {
					diamondH++
					// the initial parity is 0 since x==y at the startig position
					if (y+x)%2 != 0 {
						diamondHOdd++
					}
				}

			}
		}
	}
	// the number of hashes in the diamond shape, which are reachable by an odd number of steps
	hashesOdd := hashes - hashesEven

	/*
	   To understand this part of the solution, work on paper with a small value of `steps`, which is still of the same form as the input, like `n*131+65`. The parity of `n` is probably important, so choose n=4 or something.

	   Draw a picture of the infinite grid and the repeating diamonds inside it.

	   Note that a cell which is evenly-reachable from the centre of one grid is oddly-reachable from the centre of the grid immediately adjacent to it.

	   I will use big-diamond to refer to the large diamond mentioned in the first block comment, not the smaller diamond present in each grid square.

	   The first value subtracts any hashes which lie in a grid square which is inside the big-diamond, but for which the grid square doesn't contain the edge of the big-diamond.

	   The second value computes the hashes which lie inside the (small) diamonds, such that the small diamond shares an edge with the big-diamond.

	   The last two values compute the hashes which lie in the leftover triangles (draw a picture to see them). These are the triangles outside the small diamonds in the input.
	*/

	fullHashesForRemoval := hashesEven*(2*gridSteps-1) + hashes*(gridSteps-1)*(gridSteps-1)
	diamondHashesForRemoval := 4 * gridSteps * diamondHOdd
	finalTriangleRemoval1 := gridSteps * (hashes - diamondH)
	finalTriangleRemoval2 := (2*gridSteps - 1) * (hashesOdd - diamondHOdd)

	// final result
	fmt.Println(diamondUnobstructedReachable - fullHashesForRemoval - diamondHashesForRemoval - finalTriangleRemoval1 - finalTriangleRemoval2)
}
