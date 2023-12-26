// Why does this work?!
// Why does the heuristic not help??!
// I don't know the answers to these questions.

package main

import (
	"fmt"
	"slices"
	"strconv"
)

func compareWithHeuristic(a, b route) int {
	return compareInt(
		a.heatLoss+underestimatedCostToTarget[a.currentLocation[0]][a.currentLocation[1]],
		b.heatLoss+underestimatedCostToTarget[b.currentLocation[0]][b.currentLocation[1]],
	)
}

func compareWithoutHeuristic(a, b route) int {
	return compareInt(
		a.heatLoss,
		b.heatLoss,
	)
}

type route struct {
	currentLocation   [2]int
	path              [][2]int
	heatLoss          int
	straightLineCount int
	latestDirection   [2]int
}

var heatLossGrid [][]int
var maxX int
var maxY int

// given a route, returns a slice of all the possible next steps given the conditions:
// - no moving in a straight line for more than 3 steps in a row (at most 3 is allowed)
// - no turning back (reversing)
// - no going off the edge of the grid
func possibleNextDirections1(routeSoFar route) [][2]int {
	result := make([][2]int, 0)
	var newX, newY int
	for _, candidateDirection := range [][2]int{[2]int{1, 0}, [2]int{0, 1}, [2]int{-1, 0}, [2]int{0, -1}} {
		// condition 1
		if routeSoFar.straightLineCount == 3 && routeSoFar.latestDirection[0] == candidateDirection[0] && routeSoFar.latestDirection[1] == candidateDirection[1] {
			continue
		}

		// condition 2
		if routeSoFar.latestDirection[0] == -1*candidateDirection[0] && routeSoFar.latestDirection[1] == -1*candidateDirection[1] {
			continue
		}

		// condition 3
		newY = routeSoFar.currentLocation[0] + candidateDirection[0]
		newX = routeSoFar.currentLocation[1] + candidateDirection[1]
		if newX < 0 || newY < 0 || newX > maxX || newY > maxY {
			continue
		}

		// uncomment to imagine that doubling back on ourselves is not allowed
		//if slices.Contains(routeSoFar.path, [2]int{newY, newX}) {
		//	continue
		//}

		// all conditions have passed
		result = append(result, candidateDirection)
	}

	return result
}

func possibleNextDirections2(routeSoFar route) [][2]int {
	result := make([][2]int, 0)
	var newX, newY int
	for _, candidateDirection := range [][2]int{[2]int{1, 0}, [2]int{0, 1}, [2]int{-1, 0}, [2]int{0, -1}} {
		// condition 1
		if routeSoFar.straightLineCount == 10 && routeSoFar.latestDirection[0] == candidateDirection[0] && routeSoFar.latestDirection[1] == candidateDirection[1] {
			continue
		}

		// condition 2
		if routeSoFar.latestDirection[0] == -1*candidateDirection[0] && routeSoFar.latestDirection[1] == -1*candidateDirection[1] {
			continue
		}

		// condition 3
		newY = routeSoFar.currentLocation[0] + candidateDirection[0]
		newX = routeSoFar.currentLocation[1] + candidateDirection[1]
		if newX < 0 || newY < 0 || newX > maxX || newY > maxY {
			continue
		}

		// minimum of 4 steps before we can turn
		// (provided this is not our first move!)
		if routeSoFar.latestDirection != [2]int{0, 0} && routeSoFar.straightLineCount < 4 && (candidateDirection[0] != routeSoFar.latestDirection[0] || candidateDirection[1] != routeSoFar.latestDirection[1]) {
			continue
		}

		// all conditions have passed
		result = append(result, candidateDirection)
	}

	return result
}

// update the incomplete route with this next step
func (r *route) updateRoute(nextDirection [2]int) {
	r.currentLocation[0] = r.currentLocation[0] + nextDirection[0]
	r.currentLocation[1] = r.currentLocation[1] + nextDirection[1]

	r.path = append(r.path, r.currentLocation)

	r.heatLoss += heatLossGrid[r.currentLocation[0]][r.currentLocation[1]]

	if r.latestDirection[0] == nextDirection[0] && r.latestDirection[1] == nextDirection[1] {
		r.straightLineCount++
	} else {
		r.straightLineCount = 1
	}

	r.latestDirection = nextDirection
}

var underestimatedCostToTarget [][]int

func main() {
	lines := fetchLines(17)
	heatLossGrid = make([][]int, 0)
	for _, line := range lines {
		intLine := make([]int, len(line))
		for i, digit := range line {
			intDigit, err := strconv.Atoi(string(digit))
			check(err)
			intLine[i] = intDigit
		}
		heatLossGrid = append(heatLossGrid, intLine)
	}
	maxY = len(heatLossGrid) - 1
	maxX = len(heatLossGrid[0]) - 1

	////////////////////////////////////////////////

	// populate underestimated costs
	// fill the 2d array with zeroes
	for _ = range lines {
		estimatedCostRow := make([]int, len(lines[0]))
		underestimatedCostToTarget = append(underestimatedCostToTarget, estimatedCostRow)
	}
	// fill in the last row
	for x := maxX - 1; x >= 0; x-- {
		underestimatedCostToTarget[maxY][x] = underestimatedCostToTarget[maxY][x+1] + heatLossGrid[maxY][x+1]
	}
	// fill in all the other rows
	for y := maxY - 1; y >= 0; y-- {
		for x := 0; x <= maxX; x++ {
			underestimatedCostToTarget[y][x] = underestimatedCostToTarget[y+1][x] + heatLossGrid[y+1][x]
		}
	}
	// for debugging, you can try and inspect this:
	/*
		for _, estimateRow := range underestimatedCostToTarget {
			for _, estimateValue := range estimateRow {
				fmt.Printf("%d ", estimateValue)
			}
			fmt.Printf("\n")
		}
	*/

	// uncomment to assert that the cost everywhere but the last point is strictly greater than 0
	/*
		countZero := 0
		for _, estimateRow := range underestimatedCostToTarget {
			for _, estimateValue := range estimateRow {
				if estimateValue == 0 {
					countZero++
				}
			}
		}
		assert(countZero == 1, "estimate cost is zero in more than one place which is incorrect given our assumptions")
	*/

	// how cheaply can we get to a given location?
	locationLastDirectionAndCountToHeatLoss := make(map[[5]int]int)

	// a sorted slice of routes, in order of increasing heat loss
	incompleteRoutes := make([]route, 1)
	incompleteRoutes[0] = route{
		currentLocation:   [2]int{0, 0},
		path:              [][2]int{[2]int{0, 0}},
		heatLoss:          0,
		straightLineCount: 0,
		latestDirection:   [2]int{0, 0},
	}

	for {
		// at the beginning, assume the slice is sorted already

		// take the cheapest one incomplete route and try to extend them

		// pop from the slice
		incompleteRoute := incompleteRoutes[0]
		incompleteRoutes = incompleteRoutes[1:]

		// if we are at the destination, great!
		if incompleteRoute.currentLocation[0] == maxY && incompleteRoute.currentLocation[1] == maxX {
			//fmt.Println(incompleteRoute.path)
			fmt.Println(incompleteRoute.heatLoss)
			//fmt.Println(incompleteRoute.straightLineCount)
			//return
			break
		}
		//fmt.Printf("%d,%d %d\n", incompleteRoute.currentLocation[0], incompleteRoute.currentLocation[1], incompleteRoute.heatLoss)

		// find all the possible directions
		for _, nextDirection := range possibleNextDirections1(incompleteRoute) {
			// copy route by value
			newPath := make([][2]int, len(incompleteRoute.path))
			copy(newPath, incompleteRoute.path)
			newRoute := route{
				currentLocation:   incompleteRoute.currentLocation,
				path:              newPath,
				heatLoss:          incompleteRoute.heatLoss,
				straightLineCount: incompleteRoute.straightLineCount,
				latestDirection:   incompleteRoute.latestDirection,
			}

			// update this route with the possible next direction
			newRoute.updateRoute(nextDirection)

			// add back to the slice of incomplete routes, only if it is the cheapest way to get to the destination
			// indexToInsert, _ := slices.BinarySearchFunc(incompleteRoutes, newRoute, compareWithHeuristic)
			// incompleteRoutes = slices.Insert(incompleteRoutes, indexToInsert, newRoute)
			// /*
			newRouteIndex := [5]int{newRoute.currentLocation[0], newRoute.currentLocation[1], newRoute.latestDirection[0], newRoute.latestDirection[1], newRoute.straightLineCount}
			cheapestYet, exists := locationLastDirectionAndCountToHeatLoss[newRouteIndex]
			if !exists || newRoute.heatLoss < cheapestYet {
				indexToInsert, _ := slices.BinarySearchFunc(incompleteRoutes, newRoute, compareWithoutHeuristic)
				incompleteRoutes = slices.Insert(incompleteRoutes, indexToInsert, newRoute)
				locationLastDirectionAndCountToHeatLoss[newRouteIndex] = newRoute.heatLoss
			}
			// */

		}

		// uncomment to assert the slice is sorted
		/*
			assert(slices.IsSortedFunc(incompleteRoutes, func(a, b route) int {
				return compare(a.heatLoss+underestimatedCostToTarget[a.currentLocation[0]][a.currentLocation[1]], b.heatLoss+underestimatedCostToTarget[b.currentLocation[0]][b.currentLocation[1]])
			}), "we are not sorted")
		*/

	}

	////////////////////////////////////////////////

	// how cheaply can we get to a given location?
	locationLastDirectionAndCountToHeatLoss = make(map[[5]int]int)

	// a sorted slice of routes, in order of increasing heat loss
	incompleteRoutes = make([]route, 1)
	incompleteRoutes[0] = route{
		currentLocation:   [2]int{0, 0},
		path:              [][2]int{[2]int{0, 0}},
		heatLoss:          0,
		straightLineCount: 0,
		latestDirection:   [2]int{0, 0},
	}

	for {
		// at the beginning, assume the slice is sorted already

		// take the cheapest one incomplete route and try to extend them

		// pop from the slice
		incompleteRoute := incompleteRoutes[0]
		incompleteRoutes = incompleteRoutes[1:]

		// if we are at the destination, great!
		//fmt.Println(incompleteRoute.path)
		if incompleteRoute.currentLocation[0] == maxY && incompleteRoute.currentLocation[1] == maxX && incompleteRoute.straightLineCount >= 4 {
			fmt.Println(incompleteRoute.heatLoss)
			break
		}

		// find all the possible directions
		for _, nextDirection := range possibleNextDirections2(incompleteRoute) {
			// copy route by value
			newPath := make([][2]int, len(incompleteRoute.path))
			copy(newPath, incompleteRoute.path)
			newRoute := route{
				currentLocation:   incompleteRoute.currentLocation,
				path:              newPath,
				heatLoss:          incompleteRoute.heatLoss,
				straightLineCount: incompleteRoute.straightLineCount,
				latestDirection:   incompleteRoute.latestDirection,
			}

			// update this route with the possible next direction
			newRoute.updateRoute(nextDirection)

			// add back to the slice of incomplete routes, only if it is the cheapest way to get to the destination
			newRouteIndex := [5]int{newRoute.currentLocation[0], newRoute.currentLocation[1], newRoute.latestDirection[0], newRoute.latestDirection[1], newRoute.straightLineCount}
			cheapestYet, exists := locationLastDirectionAndCountToHeatLoss[newRouteIndex]
			if !exists || newRoute.heatLoss < cheapestYet {
				indexToInsert, _ := slices.BinarySearchFunc(incompleteRoutes, newRoute, compareWithoutHeuristic)
				incompleteRoutes = slices.Insert(incompleteRoutes, indexToInsert, newRoute)
				locationLastDirectionAndCountToHeatLoss[newRouteIndex] = newRoute.heatLoss
			}

		}

	}

}
