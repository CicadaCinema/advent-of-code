package main

import "fmt"

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

type brick struct {
	lowerCorner [3]int
	upperCorner [3]int
}

// given that there are no intersections between bricks, returns the indexes of the bricks which the target supports
func indexSupports(bricks []brick, targetIndex int) []int {
	result := make([]int, 0)

	// move this brick up
	bricks[targetIndex].lowerCorner[2]++
	bricks[targetIndex].upperCorner[2]++

	// see which indexes are now intersecting with it
	for i := range bricks {
		if i == targetIndex {
			// don't compare brick with itself
			continue
		}
		if intersect(bricks[i], bricks[targetIndex]) {
			result = append(result, i)
		}
	}

	// move the brick down again
	bricks[targetIndex].lowerCorner[2]--
	bricks[targetIndex].upperCorner[2]--

	return result
}

func indexIntersects(bricks []brick, targetIndex int) bool {
	for i := range bricks {
		if i == targetIndex {
			continue
		}
		if intersect(bricks[i], bricks[targetIndex]) {
			return true
		}
	}
	return false
}

// given two bricks, see if they intersect
// note that the order of arguments does not matter
func intersect(lower, higher brick) bool {
	if lower.lowerCorner[2] > higher.lowerCorner[2] {
		lower, higher = higher, lower
	}

	// ensure the arguments are correct
	assert(lower.lowerCorner[2] <= higher.lowerCorner[2], "arguments invalid")

	// if there is some vertical space between the bricks, clearly they do not intersect
	if lower.upperCorner[2] < higher.lowerCorner[2] {
		return false
	}

	// now we know that at least some cells in each brick take up the same z slice

	xMin, xMax := lower.lowerCorner[0], lower.upperCorner[0]
	if xMin > xMax {
		xMin, xMax = xMax, xMin
	}
	yMin, yMax := lower.lowerCorner[1], lower.upperCorner[1]
	if yMin > yMax {
		yMin, yMax = yMax, yMin
	}

	assert(xMin <= xMax, "x max assertion error")
	assert(yMin <= yMax, "y max assertion error")

	if higher.lowerCorner[0] < xMin && higher.upperCorner[0] < xMin {
		return false
	}
	if higher.lowerCorner[0] > xMax && higher.upperCorner[0] > xMax {
		return false
	}
	if higher.lowerCorner[1] < yMin && higher.upperCorner[1] < yMin {
		return false
	}
	if higher.lowerCorner[1] > yMax && higher.upperCorner[1] > yMax {
		return false
	}

	return true
}

func main() {
	lines := fetchLines(22)

	////////////////////////////////////////////////

	bricks := make([]brick, 0)

	for _, line := range lines {
		coords := integers(line)

		for _, v := range coords {
			assert(v >= 0, "negative coord value")
		}
		assert(len(coords) == 6, "invalid coord length")
		assert(coords[2] > 0, "non positive z value")
		assert(coords[5] > 0, "non positive z value")
		assert(coords[5] >= coords[2], "upper corner is not the second corner, as expected")

		bricks = append(bricks, brick{
			lowerCorner: [3]int{coords[0], coords[1], coords[2]},
			upperCorner: [3]int{coords[3], coords[4], coords[5]},
		})

	}

	// assert that no brick intersects to begin with
	for i := range bricks {
		assert(!indexIntersects(bricks, i), "two bricks intersect and we haven't even done anything")
	}

	// make the bricks fall
	for {
		movedBrick := false
		for i := range bricks {
			// if this brick is at the very bottom, we definitely can't move it
			if bricks[i].lowerCorner[2] == 1 {
				continue
			}

			// try to move this brick down and see if this causes an intersection
			bricks[i].lowerCorner[2]--
			bricks[i].upperCorner[2]--
			if !indexIntersects(bricks, i) {
				// if this does not cause an intersection, great
				movedBrick = true
				continue
			}

			// ah, this causes an intersection, so move the brick back up
			bricks[i].lowerCorner[2]++
			bricks[i].upperCorner[2]++
		}

		if !movedBrick {
			break
		}

	}

	// assert that no brick can fall without intersecting
	for i := range bricks {
		// if this brick is at the very bottom, we definitely can't move it
		if bricks[i].lowerCorner[2] == 1 {
			continue
		}
		bricks[i].lowerCorner[2]--
		bricks[i].upperCorner[2]--
		assert(indexIntersects(bricks, i), "index i does not intersects after making bricks fall")
		bricks[i].lowerCorner[2]++
		bricks[i].upperCorner[2]++
	}

	// maps from int to set<int>
	supportedByOriginal := make(map[int]map[int]bool)
	for i := range bricks {
		supportedByOriginal[i] = make(map[int]bool)
	}
	for i := range bricks {
		for _, supports := range indexSupports(bricks, i) {
			// `i` supports the brick `supports`
			supportedByOriginal[supports][i] = true
		}
	}

	sum := 0
	canBeDisintegrated := make(map[int]bool)
	// see which bricks can be disintegrated
outer:
	for i := range bricks {
		// see if i can be disintegrated, that is, it does not appear as the sole supporter of some other brick
		for j := range bricks {
			if i == j {
				continue
			}
			// if i is the sole supporter of j, then i cannot be disintegrated
			if len(supportedByOriginal[j]) == 1 && supportedByOriginal[j][i] == true {
				continue outer
			}

		}
		// this brick can be disintegrated
		canBeDisintegrated[i] = true
		sum++
	}
	fmt.Println(sum)

	////////////////////////////////////////////////

	// the bricks that are on the ground
	onGround := make(map[int]bool)
	for k, v := range supportedByOriginal {
		if len(v) == 0 {
			onGround[k] = true
		}
	}

	sum2 := 0
	for i := range bricks {
		// if this brick can be disintegrated, there is nothing to do
		if canBeDisintegrated[i] {
			continue
		}

		// clone supportedBy
		supportedBy := make(map[int]map[int]bool)
		for k, v := range supportedByOriginal {
			vNew := make(map[int]bool)
			for k_, v_ := range v {
				vNew[k_] = v_
			}
			supportedBy[k] = vNew
		}

		// imagine that brick i is going to disintegrate, and cause a chain reaction
		// these are all the bricks that will fall/change position as a result
		fall := make(map[int]bool)
		fall[i] = true

		for {
			change := false

			// remove all fallen bricks from support
			for k := range supportedBy {
				for fallenBrick := range fall {
					if supportedBy[k][fallenBrick] {
						delete(supportedBy[k], fallenBrick)
						change = true
					}
				}

				// if this change results in an empty support, and this brick wasn't on the ground to begin with, it must fall
				if !onGround[k] && len(supportedBy[k]) == 0 && !fall[k] {
					fall[k] = true
					change = true
				}
			}

			if !change {
				break
			}
		}

		add := len(fall) - 1
		sum2 += add
	}
	fmt.Println(sum2)

}
