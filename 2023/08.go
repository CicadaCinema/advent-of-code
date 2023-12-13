package main

import (
	"fmt"
	"strings"
)

type instruction struct {
	left  string
	right string
}

// https://github.com/CicadaCinema/six-small-algorithms-go/blob/d705e655f18044650dd906dab50a1703bf469d96/5/gcd_lcm.go#L8-L45
func gcd(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	if a == b {
		return a
	}

	if a < 0 {
		a *= -1
	}
	if b < 0 {
		b *= -1
	}

	var A, B int
	if a > b {
		A = a
		B = b
	} else {
		A = b
		B = a
	}

	R := A % B

	return gcd(B, R)
}

func lcm(a, b int) int {
	return (a / gcd(a, b)) * b
}

func main() {
	input := fetchInput(8)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]

	////////////////////////////////////////////////

	instructions := make(map[string]instruction, len(lines)-2)

	for _, line := range lines[2:] {
		instructions[line[0:3]] = instruction{
			left:  line[7:10],
			right: line[12:15],
		}
	}

	pos := "AAA"
	sequenceLength := len(lines[0])
	var count int
	for count = 0; pos != "ZZZ"; count++ {
		if lines[0][count%sequenceLength] == 'L' {
			pos = instructions[pos].left
		} else {
			pos = instructions[pos].right
		}
	}
	fmt.Println(count)

	////////////////////////////////////////////////

	startingPositions := make([]string, 0)
	for i := range instructions {
		if i[2] == 'A' {
			startingPositions = append(startingPositions, i)
		}
	}

	lcmInputs := make([]int, 0)
outer:
	for _, pos := range startingPositions {
		// start with the position defined in the outer loop,
		// and iterate until we hit a finishing position (something that ends in Z)
		// assume:
		// 1) for each starting position, there is only one finishing position (this is probably not true in general, but the input is 'nice')
		// 2) following the instructions makes us enter a cycle (this is probably true in general)
		// 3) the period of the cycle is exactly equal to the number of steps it took to get to the finishing position
		// 4) this period is also congruent to 0, mod sequenceLength
		// the last two assumptions are baffling but they seem to hold...
		for count = 0; true; count++ {
			if pos[2] == 'Z' {
				lcmInputs = append(lcmInputs, count)
				// test assumption 4) by uncommenting this line
				// fmt.Printf("pos: %s count: %d countMODsequenceLength: %d\n", pos, count, count%sequenceLength)
				// in addition, test assumption 3) by commenting this line out and inspecting the first few lines of debug output (uncommented above)
				continue outer
			}
			if lines[0][count%sequenceLength] == 'L' {
				pos = instructions[pos].left
			} else {
				pos = instructions[pos].right
			}
		}

	}

	lcmAll := lcmInputs[0]

	for _, v := range lcmInputs[1:] {
		lcmAll = lcm(lcmAll, v)
	}

	fmt.Println(lcmAll)

}
