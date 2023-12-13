package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	input := fetchInput(6)

	////////////////////////////////////////////////

	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]
	if len(lines) != 2 {
		panic("assertion failed: wrong number of lines")
	}
	rNum := regexp.MustCompile("[0-9]+")
	timeDistanceTuples := make([][2]int, 0)
	times := rNum.FindAllString(lines[0], -1)
	distances := rNum.FindAllString(lines[1], -1)
	for i := range times {
		time, err := strconv.Atoi(times[i])
		check(err)
		distance, err := strconv.Atoi(distances[i])
		check(err)
		timeDistanceTuples = append(timeDistanceTuples, [2]int{time, distance})
	}

	product := 1
	for _, tuple := range timeDistanceTuples {
		waysToWinHere := 0
		for i := 0; i <= tuple[0]; i++ {
			if i*(tuple[0]-i) > tuple[1] {
				waysToWinHere += 1
			}
		}
		product *= waysToWinHere
	}

	fmt.Println(product)

	////////////////////////////////////////////////

	rSpace := regexp.MustCompile(" ")
	bigTime := rSpace.ReplaceAllString(lines[0], "")
	bigDistance := rSpace.ReplaceAllString(lines[1], "")
	timeInt, err := strconv.Atoi(rNum.FindString(bigTime))
	time := float64(timeInt)
	check(err)
	distanceInt, err := strconv.Atoi(rNum.FindString(bigDistance))
	distance := float64(distanceInt)
	check(err)

	criticalMinus := (-time - math.Sqrt(math.Pow(time, 2)-4*distance)) / -2
	criticalPlus := (-time + math.Sqrt(math.Pow(time, 2)-4*distance)) / -2

	// inspected manually to see which is the lower bound and which is the upper bound

	criticalPlus = math.Ceil(criticalPlus)
	criticalMinus = math.Floor(criticalMinus)

	fmt.Println(int(criticalMinus) - int(criticalPlus) + 1)
}
