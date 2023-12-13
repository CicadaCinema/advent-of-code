package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sum int
var gears map[[2]int][2]int

func addToSum(partNumber string) {
	num, err := strconv.Atoi(partNumber)
	check(err)
	sum += num
}

func addToGears(partNumber string, key [2]int) {
	num, err := strconv.Atoi(partNumber)
	check(err)
	prevValue, present := gears[key]
	if !present {
		gears[key] = [2]int{num, 0}
	} else {
		prevValue[1] = num
		gears[key] = prevValue
	}
}

func main() {
	input := fetchInput(3)

	////////////////////////////////////////////////

	rNum := regexp.MustCompile("[0-9]+")
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]
	minIndex := 0
	maxIndex := len(lines[0]) - 1
	sum = 0

	for i, line := range lines {
	numbers1:
		for _, numMatch := range rNum.FindAllStringIndex(line, -1) {
			startHorizontal := minIndex
			endHorizontal := maxIndex
			if numMatch[0]-1 > startHorizontal {
				startHorizontal = numMatch[0] - 1
			}
			if numMatch[1] < endHorizontal {
				endHorizontal = numMatch[1]
			}

			if i > 0 {
				for j := startHorizontal; j <= endHorizontal; j++ {
					if lines[i-1][j] != '.' {
						addToSum(line[numMatch[0]:numMatch[1]])
						continue numbers1
					}
				}
			}
			if i < len(lines)-1 {
				for j := startHorizontal; j <= endHorizontal; j++ {
					if lines[i+1][j] != '.' {
						addToSum(line[numMatch[0]:numMatch[1]])
						continue numbers1
					}
				}
			}

			if (numMatch[0] > 0 && line[numMatch[0]-1] != '.') || (numMatch[1] < len(line) && line[numMatch[1]] != '.') {
				addToSum(line[numMatch[0]:numMatch[1]])
				continue numbers1
			}

		}
	}

	fmt.Println(sum)

	////////////////////////////////////////////////

	sum = 0
	gears = make(map[[2]int][2]int)

	for i, line := range lines {
	numbers2:
		for _, numMatch := range rNum.FindAllStringIndex(line, -1) {
			startHorizontal := minIndex
			endHorizontal := maxIndex
			if numMatch[0]-1 > startHorizontal {
				startHorizontal = numMatch[0] - 1
			}
			if numMatch[1] < endHorizontal {
				endHorizontal = numMatch[1]
			}

			if i > 0 {
				for j := startHorizontal; j <= endHorizontal; j++ {
					if lines[i-1][j] == '*' {
						addToGears(line[numMatch[0]:numMatch[1]], [2]int{i - 1, j})
						continue numbers2
					}
				}
			}
			if i < len(lines)-1 {
				for j := startHorizontal; j <= endHorizontal; j++ {
					if lines[i+1][j] == '*' {
						addToGears(line[numMatch[0]:numMatch[1]], [2]int{i + 1, j})
						continue numbers2
					}
				}
			}

			if numMatch[0] > 0 && line[numMatch[0]-1] == '*' {
				addToGears(line[numMatch[0]:numMatch[1]], [2]int{i, numMatch[0] - 1})
				continue numbers2
			}

			if numMatch[1] < len(line) && line[numMatch[1]] == '*' {
				addToGears(line[numMatch[0]:numMatch[1]], [2]int{i, numMatch[1]})
				continue numbers2
			}
		}
	}

	for _, v := range gears {
		sum += v[0] * v[1]
	}

	fmt.Println(sum)

}
