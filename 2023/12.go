package main

import (
	"fmt"
	"strings"
)

var possibilitiesCache map[[2]int]int

// the number of possibilities of assigning blocks in this line, provided it is legal to begin a block directly at index 0
func numPossibilities(line string, blockLengths []int, blockLengthSum int) int {
	if blockLengthSum+len(blockLengths)-1 > len(line) {
		// not enough space
		return 0
	}

	cacheKey := [2]int{len(line), blockLengthSum}
	result, cacheHit := possibilitiesCache[cacheKey]
	if cacheHit {
		return result
	}

	pos := 0

	// base case - we only need to assign one block and then we are done
	if len(blockLengths) == 1 {
		lastHashIndex := -1
		for i, v := range line {
			if v == '#' {
				lastHashIndex = i
			}
		}
		for startIndex := 0; startIndex <= len(line)-blockLengths[0]; startIndex++ {
			// want:
			// -- no hash -- | -- no dot -- | -- no hash --
			if startIndex > 0 && line[startIndex-1] == '#' {
				break
			}
			if lastHashIndex >= startIndex+blockLengths[0] {
				continue
			}
			if !strings.ContainsRune(line[startIndex:startIndex+blockLengths[0]], '.') {
				pos++
			}
		}
	} else {
		// the upper bound on startIndex does not currently take into account the gaps in between the remaining blocks and so it can be more strict if necessary
		// but that condition is sort of already handled at the top of this function so it is no big deal really
		for startIndex := 0; startIndex <= len(line)-blockLengthSum; startIndex++ {
			// want:
			// -- no hash -- | -- no dot -- | rune-which-is-not-hash
			if startIndex > 0 && line[startIndex-1] == '#' {
				break
			}
			if !strings.ContainsRune(line[startIndex:startIndex+blockLengths[0]], '.') && line[startIndex+blockLengths[0]] != '#' {
				pos += numPossibilities(line[startIndex+blockLengths[0]+1:], blockLengths[1:], blockLengthSum-blockLengths[0])
			}
		}
	}

	_, cacheHit = possibilitiesCache[cacheKey]
	if cacheHit {
		panic("failed assertion, expected to place new value into the cache")
	}
	possibilitiesCache[cacheKey] = pos
	return pos
}

func main() {
	lines := fetchLines(12)

	sum1 := 0
	sum2 := 0

	for _, line := range lines {
		parts := strings.Split(line, " ")

		if len(parts) != 2 {
			panic("invalid input data on this line")
		}

		possibilitiesCache = make(map[[2]int]int)
		lengths := uintegers(parts[1])
		lengthSum := sliceSum(lengths)
		sum1 += numPossibilities(parts[0], lengths, lengthSum)

		possibilitiesCache = make(map[[2]int]int)
		lengths = uintegers(strings.Repeat(","+parts[1], 5)[1:])
		sum2 += numPossibilities(strings.Repeat("?"+parts[0], 5)[1:], lengths, lengthSum*5)
	}

	fmt.Println(sum1)
	fmt.Println(sum2)
}
