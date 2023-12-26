package main

import (
	"fmt"
	"strconv"
	"strings"
)

func hash(inputString string) int {
	result := 0
	for _, v := range inputString {
		result += int(v)
		result *= 17
		result = result % 256
	}
	return result
}

type lens struct {
	label       string
	focalLength int
}

var state [256][]lens

func processOperation(op string) [256][]lens {
	equalsSplit := strings.Split(op, "=")
	if len(equalsSplit) == 1 {
		// this is a - instruction
		label := op[:len(op)-1]
		boxIndex := hash(label)
		newLenses := make([]lens, 0)
		for _, lens := range state[boxIndex] {
			if lens.label != label {
				newLenses = append(newLenses, lens)
			}
		}
		state[boxIndex] = newLenses
	} else {
		// this is a = instruction
		focalLength, err := strconv.Atoi(equalsSplit[1])
		check(err)

		boxIndex := hash(equalsSplit[0])

		// check if there is already a lens with this label
		for i, lens := range state[boxIndex] {
			if lens.label == equalsSplit[0] {
				state[boxIndex][i].focalLength = focalLength
				return state
			}
		}

		state[boxIndex] = append(state[boxIndex], lens{
			label:       equalsSplit[0],
			focalLength: focalLength,
		})

	}
	return state
}

func main() {
	lines := fetchLines(15)
	assert(len(lines) == 1, "only expected one line, got more")

	state := [256][]lens{}
	for i := range state {
		state[i] = make([]lens, 0)
	}

	sum := 0
	for _, element := range strings.Split(lines[0], ",") {
		sum += hash(element)
		state = processOperation(element)
	}
	fmt.Println(sum)

	// add up the focusing power of all the lenses
	sum = 0
	for i := range state {
		for j, lens := range state[i] {
			sum += (1 + i) * (1 + j) * lens.focalLength
		}
	}
	fmt.Println(sum)

}
