package main

import "fmt"

type pos struct {
	x int
	y int
}

var rowIsEmpty []bool
var colIsEmpty []bool

var galaxies []pos

func distanceBetweenPair(a, b pos) (int, int) {
	// assume that b.y>=a.y
	if !(b.y >= a.y) {
		panic("failed assetion in determining distances between galaxies")
	}

	result1 := b.y - a.y
	result2 := b.y - a.y
	// row expansion
	for rowIndex := a.y + 1; rowIndex < b.y; rowIndex++ {
		if rowIsEmpty[rowIndex] {
			result1 += 1
			result2 += 999999
		}
	}

	// col expansion
	if b.x >= a.x {
		result1 += b.x - a.x
		result2 += b.x - a.x
		for colIndex := a.x + 1; colIndex < b.x; colIndex++ {
			if colIsEmpty[colIndex] {
				result1 += 1
				result2 += 999999
			}
		}
	} else {
		result1 += a.x - b.x
		result2 += a.x - b.x
		for colIndex := b.x + 1; colIndex < a.x; colIndex++ {
			if colIsEmpty[colIndex] {
				result1 += 1
				result2 += 999999
			}
		}
	}

	return result1, result2
}

func main() {
	lines := fetchLines(11)

	// is the row/col with the given index empty of galaxies?
	rowIsEmpty = make([]bool, len(lines))
	colIsEmpty = make([]bool, len(lines[0]))

	// initially assume everything is empty
	for i := range rowIsEmpty {
		rowIsEmpty[i] = true
	}
	for i := range colIsEmpty {
		colIsEmpty[i] = true
	}

	galaxies = make([]pos, 0)

	for y := range lines {
		for x := range lines[y] {
			if lines[y][x] == '#' {
				galaxies = append(galaxies, pos{x: x, y: y})
				rowIsEmpty[y] = false
				colIsEmpty[x] = false
			}
		}
	}

	// assume that the number of galaxies is sufficiently large
	sum1 := 0
	sum2 := 0
	for i := range galaxies {
		for j := i + 1; j < len(galaxies); j++ {
			d1, d2 := distanceBetweenPair(galaxies[i], galaxies[j])
			sum1 += d1
			sum2 += d2
		}
	}
	fmt.Println(sum1)
	fmt.Println(sum2)
}
