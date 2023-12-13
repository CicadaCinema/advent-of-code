package main

import (
	"fmt"
	"strconv"
	"strings"

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

var source2dest [][][3]int

func seed2location(seed int) int {
	result := seed

	// do let(source2dest) transformations of the result
	for _, transformation := range source2dest {
		ruleIndex, found := slices.BinarySearchFunc(transformation, result, func(tuple [3]int, target int) int {
			return compare(tuple[1], target)
		})

		if found {
			// we have found an exact match for one of the rules
			result = transformation[ruleIndex][0]
			continue
		}

		if ruleIndex == 0 {
			// the element target appears before any rule
			continue
		}

		closestRule := transformation[ruleIndex-1]

		// if the range is large enough, apply the rule
		// otherwise do nothing
		if result-closestRule[1] < closestRule[2] {
			if result < closestRule[1] {
				panic("assertion failed: the result somehow preceeds the chosen rule")
			}
			result = result + closestRule[0] - closestRule[1]
		}
	}

	return result
}

func seedRange2smallestLocation(seedRange []int) int {
	if len(seedRange) != 2 {
		panic("assertion failed: we expect the seed range to contain exactly two elements")
	}
	lowestLocation := seed2location(seedRange[0])
	for i := seedRange[0] + 1; i < seedRange[0]+seedRange[1]; i++ {
		thisLocation := seed2location(i)
		if thisLocation < lowestLocation {
			lowestLocation = thisLocation
		}
	}
	return lowestLocation
}

func main() {
	lines := strings.Split(fetchInput(5), "\n")
	lines = lines[0 : len(lines)-1]

	// parse the target seeds
	stringSeeds := strings.Split(lines[0], " ")[1:]
	seeds := make([]int, len(stringSeeds))
	for i, stringSeed := range stringSeeds {
		seed, err := strconv.Atoi(stringSeed)
		check(err)
		seeds[i] = seed
	}

	// a slice containing slices containing 3-element source-to-dest tuples
	source2dest = make([][][3]int, 0)
	currentS2DIndex := -1

	for _, line := range lines[2:] {
		// ignore empty lines
		if line == "" {
			continue
		}

		// lines starting with a letter make us move to the next s2d slice
		if int(line[0]) > 57 {
			currentS2DIndex += 1
			source2dest = append(source2dest, make([][3]int, 0))
			continue
		}

		// now we know we are on a line with a 3-tuple
		stringTuple := strings.Split(line, " ")
		if len(stringTuple) != 3 {
			panic("tuple has unexpected length")
		}

		tuple := [3]int{}

		for i, v := range stringTuple {
			element, err := strconv.Atoi(v)
			check(err)
			tuple[i] = element
		}

		source2dest[currentS2DIndex] = append(source2dest[currentS2DIndex], tuple)
	}

	// sort all the s2d slices
	for _, s2dSlice := range source2dest {
		slices.SortFunc(s2dSlice, func(a, b [3]int) int {
			return compare(a[1], b[1])
		})
	}

	smallestLocation := seed2location(seeds[0])
	for i := 1; i < len(seeds); i++ {
		thisLocation := seed2location(seeds[i])
		if thisLocation < smallestLocation {
			smallestLocation = thisLocation
		}
	}

	fmt.Println(smallestLocation)

	smallestLocation = seedRange2smallestLocation(seeds[0:2])
	for i := 2; i < len(seeds); i += 2 {
		thisLocation := seedRange2smallestLocation(seeds[i : i+2])
		if thisLocation < smallestLocation {
			smallestLocation = thisLocation
		}
	}

	fmt.Println(smallestLocation)
}
