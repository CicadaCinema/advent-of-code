package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func compare(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return +1
	}
	return 0
}

func main() {
	input := fetchInput(7)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]

	if len(lines) != 1000 {
		panic("invalid number of lines in input")
	}

	cardStrengths := make(map[byte]int)
	cardStrengths['A'] = 12
	cardStrengths['K'] = 11
	cardStrengths['Q'] = 10
	cardStrengths['J'] = 9
	cardStrengths['T'] = 8
	cardStrengths['9'] = 7
	cardStrengths['8'] = 6
	cardStrengths['7'] = 5
	cardStrengths['6'] = 4
	cardStrengths['5'] = 3
	cardStrengths['4'] = 2
	cardStrengths['3'] = 1
	cardStrengths['2'] = 0

	cardStrengthsJoker := make(map[byte]int)
	cardStrengthsJoker['A'] = 12
	cardStrengthsJoker['K'] = 11
	cardStrengthsJoker['Q'] = 10
	cardStrengthsJoker['T'] = 9
	cardStrengthsJoker['9'] = 8
	cardStrengthsJoker['8'] = 7
	cardStrengthsJoker['7'] = 6
	cardStrengthsJoker['6'] = 5
	cardStrengthsJoker['5'] = 4
	cardStrengthsJoker['4'] = 3
	cardStrengthsJoker['3'] = 2
	cardStrengthsJoker['2'] = 1
	cardStrengthsJoker['J'] = 0

	// index 1 - strength for part 1
	// index 2 - strength for part 2
	bidAndStrengthTuples := make([][3]int, len(lines))

	for i := range lines {
		hand := lines[i][:5]
		bid, err := strconv.Atoi(lines[i][6:])
		check(err)
		if len(hand) != 5 {
			panic("card has unexpected length")
		}

		var typeStrength int
		occurences := make(map[byte]int)
		for _, handCard := range hand {
			occurences[byte(handCard)] += 1
		}
		sortedFrequencies := maps.Values(occurences)
		sort.Sort(sort.Reverse(sort.IntSlice(sortedFrequencies)))
		if sortedFrequencies[0] == 5 {
			typeStrength = 7
		} else if sortedFrequencies[0] == 4 {
			typeStrength = 6
		} else if sortedFrequencies[0] == 3 {
			if sortedFrequencies[1] == 2 {
				typeStrength = 5
			} else {
				typeStrength = 4
			}
		} else if sortedFrequencies[0] == 2 {
			if sortedFrequencies[1] == 2 {
				typeStrength = 3
			} else {
				typeStrength = 2
			}
		} else {
			typeStrength = 1
		}

		bidAndStrengthTuples[i] = [3]int{bid, typeStrength*1000000 + cardStrengths[hand[0]]*13*13*13*13 + cardStrengths[hand[1]]*13*13*13 + cardStrengths[hand[2]]*13*13 + cardStrengths[hand[3]]*13 + cardStrengths[hand[4]], 0}
		//fmt.Printf("line: %s || bid: %d || hand: %s || strength: %d\n", lines[i], bid, hand, bidAndStrengthTuples[i][1])

		// part 2
		switch occurences['J'] {
		case 5:
			// 5 jokers, can make 5 of a kind
			typeStrength = 7
		case 4:
			// 4 jokers, one of another type, can make 5 of a kind
			typeStrength = 7
		case 3:
			// 3 jokers, the second-most-common card can appear either 2 or 1 times
			if sortedFrequencies[1] == 2 {
				// 3 jokers, 2x another type, can make 5 of a kind
				typeStrength = 7
			} else {
				// 3 jokers, 1x another, 1x another, can make 4 of a kind
				typeStrength = 6
			}
		case 2:
			if sortedFrequencies[0] == 3 {
				// 2 jokers, 3x another, can make 5 of a kind
				typeStrength = 7
			} else {
				// two jokers, every card appears either twice or one time
				if sortedFrequencies[1] == 2 {
					// two jokers, 2x another, 1x another, can make 4 of a kind
					typeStrength = 6
				} else {
					// two jokers, 1x another, 1x another, 1x another, can make 3 of a kind
					typeStrength = 4
				}
			}
		case 1:
			if sortedFrequencies[0] == 4 {
				// 1 joker, 4x another, can make 5 of a kind
				typeStrength = 7
			} else if sortedFrequencies[0] == 3 {
				// 1 joker, 3x another, can make 4 of a kind
				typeStrength = 6
			} else if sortedFrequencies[0] == 2 {
				if sortedFrequencies[1] == 2 {
					// 1 joker, 2x another, 2x another, can make full house
					typeStrength = 5
				} else {
					// 1 joker, 2x another, 1x another, 1x another, can make 3 of a kind
					typeStrength = 4
				}
			} else {
				// every card appears once, can make one pair
				typeStrength = 2
			}
		case 0:
			// do nothing
		default:
			panic("invalid number of occurences of J")
		}

		bidAndStrengthTuples[i][2] = typeStrength*1000000 + cardStrengthsJoker[hand[0]]*13*13*13*13 + cardStrengthsJoker[hand[1]]*13*13*13 + cardStrengthsJoker[hand[2]]*13*13 + cardStrengthsJoker[hand[3]]*13 + cardStrengthsJoker[hand[4]]
	}

	// sort in ascending order of strength for part 1
	slices.SortFunc(bidAndStrengthTuples, func(a, b [3]int) int {
		return compare(a[1], b[1])
	})

	sum := 0
	for i := range bidAndStrengthTuples {
		//fmt.Println(bidAndStrengthTuples[i])
		sum += bidAndStrengthTuples[i][0] * (i + 1)
	}
	fmt.Println(sum)

	// sort in ascending order of strength for part 2
	slices.SortFunc(bidAndStrengthTuples, func(a, b [3]int) int {
		return compare(a[2], b[2])
	})

	sum = 0
	for i := range bidAndStrengthTuples {
		//fmt.Println(bidAndStrengthTuples[i])
		sum += bidAndStrengthTuples[i][0] * (i + 1)
	}
	fmt.Println(sum)
}
