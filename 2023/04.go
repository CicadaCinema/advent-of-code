package main

import (
	"fmt"
	"strings"
)

func main() {
	input := fetchInput(4)

	////////////////////////////////////////////////

	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]
	sum := 0

	var cardWinnings int
	var winningNums map[string]bool

	for _, line := range lines {
		cardWinnings = 0
		winningNums = make(map[string]bool)

		// remove leading "Card xyz: " characters
		line = line[10:len(line)]

		winningAndActual := strings.Split(line, " | ")

		// populate winning nums
		for _, winningString := range strings.Split(winningAndActual[0], " ") {
			// account for adjacent separators
			if winningString == "" {
				continue
			}
			winningNums[winningString] = true
		}

		// count up the winnings
		for _, actualString := range strings.Split(winningAndActual[1], " ") {
			// account for adjacent separators
			if actualString == "" {
				continue
			}
			if winningNums[actualString] {
				if cardWinnings == 0 {
					cardWinnings = 1
				} else {
					cardWinnings *= 2
				}
			}
		}

		sum += cardWinnings
	}

	fmt.Println(sum)

	////////////////////////////////////////////////

	// the number of scratchcards of each type that you own
	scratchcardCounts := make([]int, len(lines))
	for i := range scratchcardCounts {
		scratchcardCounts[i] = 1
	}

	for scratchcardIndex, line := range lines {
		cardWinnings = 0
		winningNums = make(map[string]bool)

		// remove leading "Card xyz: " characters
		line = line[10:len(line)]

		winningAndActual := strings.Split(line, " | ")

		// populate winning nums
		for _, winningString := range strings.Split(winningAndActual[0], " ") {
			// account for adjacent separators
			if winningString == "" {
				continue
			}
			winningNums[winningString] = true
		}

		// count up the winnings
		for _, actualString := range strings.Split(winningAndActual[1], " ") {
			// account for adjacent separators
			if actualString == "" {
				continue
			}
			if winningNums[actualString] {
				cardWinnings += 1
			}
		}

		// I own this card scratchcardCounts[scratchcardIndex] times
		// each copy of it gives me the next cardWinnings cards

		for i := scratchcardIndex + 1; i <= scratchcardIndex+cardWinnings; i++ {
			scratchcardCounts[i] += scratchcardCounts[scratchcardIndex]
		}
	}

	sum = 0
	for _, v := range scratchcardCounts {
		sum += v
	}
	fmt.Println(sum)
}
