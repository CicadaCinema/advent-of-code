package main

import (
	"fmt"
)

func allZero(seq []int) bool {
	for _, v := range seq {
		if v != 0 {
			return false
		}
	}
	return true
}

func generatePrediction(data []int) (int, int) {
	// we know from inspecting the data that each data sequence will never be all zero initially
	differences := make([][]int, 1)
	differences[0] = data

	// generate more differences until the last slice is all zero
	for !allZero(differences[len(differences)-1]) {
		lastDifference := differences[len(differences)-1]
		// the length of the next differences slice is one fewer than the last difference slice
		nextDifference := make([]int, len(lastDifference)-1)
		// populate it
		for i := range nextDifference {
			nextDifference[i] = lastDifference[i+1] - lastDifference[i]
		}
		differences = append(differences, nextDifference)
	}

	// now we know that the last difference slice is all zero
	// want to get the original input data
	extrapolatedForwards := 0
	extrapolatedBackwards := 0
	for i := len(differences) - 2; i >= 0; i-- {
		difference := differences[i]

		lastGuy := difference[len(difference)-1]
		extrapolatedForwards = lastGuy + extrapolatedForwards

		extrapolatedBackwards = difference[0] - extrapolatedBackwards
	}

	return extrapolatedForwards, extrapolatedBackwards
}

func main() {
	lines := fetchLines(9)

	sumForwards := 0
	sumBackwards := 0
	for _, line := range lines {
		f, b := generatePrediction(integers(line))
		sumForwards += f
		sumBackwards += b
	}
	fmt.Println(sumForwards)
	fmt.Println(sumBackwards)
}
