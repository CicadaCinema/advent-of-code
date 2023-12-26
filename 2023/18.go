package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type corner struct {
	x int
	y int
	// "L" or "R"
	kind string
}

type boundary struct {
	// horizontal trenches, not including their corners
	// maps from vertical index (y) to inclusive range of horizontal indexes (x)
	horizontals map[int][][2]int
	// vertical trenches, not including their corners
	// maps from horizontal index (x) to inclusive range of vertical indexes (y)
	verticals   map[int][][2]int
	corners     []corner
	borderCount int
}

type originalDirection struct {
	bearing string
	length  int
}

// returns the original direction for part 1 and part 2
func lineToOriginalDirections(line string) (originalDirection, originalDirection) {
	elements := strings.Split(line, " ")
	assert(len(elements) == 3, "elements has wrong length")

	colourCode := elements[2][2 : len(elements[2])-1]
	assert(len(colourCode) == 6, "colour code has wrong length, expected 6")

	length1, err := strconv.Atoi(elements[1])
	check(err)

	length2_64, err := strconv.ParseInt(colourCode[:5], 16, 0)
	length2 := int(length2_64)

	_, exists := bearings[elements[0]]
	assert(exists, "invalid direction")
	direction1 := elements[0]

	var direction2 string
	switch colourCode[5] {
	case '3':
		direction2 = "U"
	case '2':
		direction2 = "L"
	case '1':
		direction2 = "D"
	case '0':
		direction2 = "R"
	default:
		panic("unexpected direction")
	}

	return originalDirection{
			bearing: direction1,
			length:  length1,
		}, originalDirection{
			bearing: direction2,
			length:  length2,
		}
}

func originalToBoundary(directions []originalDirection) boundary {
	result := boundary{
		horizontals: map[int][][2]int{},
		verticals:   map[int][][2]int{},
		corners:     []corner{},
	}

	x, y := 0, 0

	// count the number of steps on the boundary
	boundaryCount1 := 0

	for i := range directions {
		thisIndex, exists1 := bearings[directions[i].bearing]
		prevIndex, exists2 := bearings[directions[(i-1+len(directions))%len(directions)].bearing]
		assert(exists1, "could not find index 1")
		assert(exists2, "could not find index 2")
		switch thisIndex - prevIndex {
		case -3:
			fallthrough
		case 1:
			// right turn
			result.corners = append(result.corners, corner{
				x:    x,
				y:    y,
				kind: "R",
			})
		case 3:
			fallthrough
		case -1:
			// left turn
			result.corners = append(result.corners, corner{
				x:    x,
				y:    y,
				kind: "L",
			})
		default:
			panic("invalid turn")
		}

		var directionArray [2]int
		switch directions[i].bearing {
		case "U":
			if directions[i].length > 1 {
				_, e := result.verticals[x]
				if !e {
					result.verticals[x] = make([][2]int, 0)
				}
				result.verticals[x] = append(result.verticals[x], [2]int{y - directions[i].length + 1, y - 1})
			}
			directionArray = [2]int{-1, 0}
		case "D":
			if directions[i].length > 1 {
				_, e := result.verticals[x]
				if !e {
					result.verticals[x] = make([][2]int, 0)
				}
				result.verticals[x] = append(result.verticals[x], [2]int{y + 1, y + directions[i].length - 1})
			}
			directionArray = [2]int{1, 0}
		case "L":
			if directions[i].length > 1 {
				_, e := result.horizontals[y]
				if !e {
					result.horizontals[y] = make([][2]int, 0)
				}
				result.horizontals[y] = append(result.horizontals[y], [2]int{x - directions[i].length + 1, x - 1})
			}
			directionArray = [2]int{0, -1}
		case "R":
			if directions[i].length > 1 {
				_, e := result.horizontals[y]
				if !e {
					result.horizontals[y] = make([][2]int, 0)
				}
				result.horizontals[y] = append(result.horizontals[y], [2]int{x + 1, x + directions[i].length - 1})
			}
			directionArray = [2]int{0, 1}
		default:
			panic("unexpected direction 1")
		}

		x += directionArray[1] * directions[i].length
		y += directionArray[0] * directions[i].length

		boundaryCount1 += directions[i].length
	}

	// sort horizontals and verticals
	for k := range result.horizontals {
		slices.SortFunc(result.horizontals[k], func(a, b [2]int) int {
			return compareInt(a[0], b[0])
		})
	}
	for k := range result.verticals {
		slices.SortFunc(result.verticals[k], func(a, b [2]int) int {
			return compareInt(a[0], b[0])
		})
	}

	// sort corners by their x value
	slices.SortFunc(result.corners, func(a, b corner) int {
		return compareInt(a.x, b.x)
	})

	// assertions
	boundaryCount2 := len(result.corners)

	for _, v := range result.horizontals {
		lower := 0
		for i, horizontal := range v {
			assert(horizontal[0] <= horizontal[1], "bad horizontal")
			boundaryCount2 += horizontal[1] - horizontal[0] + 1
			if i != 0 {
				assert(lower < horizontal[0], "bad sorting")
			}
			lower = horizontal[0]
		}
	}
	for _, v := range result.verticals {
		lower := 0
		for i, vertical := range v {
			assert(vertical[0] <= vertical[1], "bad vertical")
			boundaryCount2 += vertical[1] - vertical[0] + 1
			if i != 0 {
				assert(lower < vertical[0], "bad sorting")
			}
			lower = vertical[0]
		}
	}

	assert(boundaryCount1 == boundaryCount2, "bad boundary count")

	result.borderCount = boundaryCount1

	return result
}

type indicator struct {
	// "L", "R" or "|"
	kind  string
	index int
	// whether or not this indicator causes a flip in inside-ness (for example, always true for indicators of kind "|")
	flip bool
	// if this flip results in inside-ness
	insideRegionBegin bool
}

func countOnLine(b boundary, y int) int {
	line := []indicator{}

	// populate verticals
	for k, vSlice := range b.verticals {
		for _, v := range vSlice {
			if v[0] <= y && y <= v[1] {
				line = append(line, indicator{kind: "|", index: k})
			}
		}
	}

	// add all corners
	for _, corner := range b.corners {
		if corner.y == y {
			line = append(line, indicator{kind: corner.kind, index: corner.x})
		}
	}

	// sort indicators by index
	slices.SortFunc(line, func(a, b indicator) int {
		return compareInt(a.index, b.index)
	})

	cornerCount := 0
	for i := range line {
		if line[i].kind == "|" {
			line[i].flip = true
			continue
		}
		cornerCount++

		// now we know the kind of this indicator is either "L" or "R"
		// if the previous indicator was also a corner, like this one, and it was a different corner, then it's time to flip starting with this indicator
		if cornerCount%2 == 0 && line[i-1].kind != "|" && line[i-1].kind != line[i].kind {
			line[i].flip = true
		}
	}

	flipCount := 0
	cornerCount = 0
	for i := range line {
		if line[i].flip {
			flipCount++
			if flipCount%2 == 1 {
				line[i].insideRegionBegin = true
			}
		}

		if line[i].kind != "|" {
			cornerCount++
			if cornerCount%2 == 0 && i >= 2 && line[i].kind == line[i-1].kind && line[i-2].insideRegionBegin {
				line[i].insideRegionBegin = true
			}
		}
	}

	sum := 0
	for i := range line {
		if line[i].insideRegionBegin {
			sum += line[i+1].index - 1 - line[i].index
		}
	}

	return sum
}

// count the number of cells inside the boundary
func countFast(b boundary) int {

	count := b.borderCount

	// interesting y are those with corners
	interestingYMap := make(map[int]bool)
	interestingY := make([]int, 0)

	for _, corner := range b.corners {
		interestingYMap[corner.y] = true
	}

	for k := range interestingYMap {
		interestingY = append(interestingY, k)
	}

	slices.Sort(interestingY)

	for i, y := range interestingY {
		// there are interesting rows with corners ...
		count += countOnLine(b, y)

		if i == 0 {
			continue
		}

		nextRepeat := y - interestingY[i-1] - 1

		if nextRepeat == 0 {
			continue
		}

		// ... and without corners (in between those interesting rows, everything repeats row by row)
		count += countOnLine(b, y-1) * nextRepeat
	}

	return count
}

// count the number of cells inside the boundary
func countSlow(b boundary) int {
	count := b.borderCount

	for y := -1000; y <= 1000; y++ {
		inside := false
		prevCorner := ""

	outer:
		for x := -1000; x <= 1000; x++ {
			// if we are in a vertical
			for _, v := range b.verticals[x] {
				if v[0] <= y && y <= v[1] {
					inside = !inside
					continue outer
				}
			}

			// if we are in a horizontal
			for _, h := range b.horizontals[y] {
				if h[0] <= x && x <= h[1] {
					continue outer
				}
			}

			// if we are at a corner
			for _, c := range b.corners {
				if c.x == x && c.y == y {
					if prevCorner == "" {
						prevCorner = c.kind
					} else if prevCorner != c.kind {
						inside = !inside
					}
					continue outer
				}
			}

			// otherwise
			if inside {
				count++
			}
		}
	}

	return count

}

var bearings map[string]int

func main() {
	lines := fetchLines(18)

	bearings = make(map[string]int)
	bearings["U"] = 0
	bearings["R"] = 1
	bearings["D"] = 2
	bearings["L"] = 3

	// populate original directions
	originalDirections1 := make([]originalDirection, len(lines))
	originalDirections2 := make([]originalDirection, len(lines))
	for i := range lines {
		originalDirections1[i], originalDirections2[i] = lineToOriginalDirections(lines[i])
	}

	boundary1 := originalToBoundary(originalDirections1)
	boundary2 := originalToBoundary(originalDirections2)

	// fmt.Println(countSlow(boundary))
	fmt.Println(countFast(boundary1))
	fmt.Println(countFast(boundary2))

}
