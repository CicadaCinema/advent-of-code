package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func tiltNorthAndCalculateLoadInColumn(column string) int {
	load := 0
	freeIndex := len(column)

	for i, v := range column {
		if v == 'O' {
			load += freeIndex
			freeIndex--
		} else if v == '#' {
			freeIndex = len(column) - i - 1
		}
	}

	return load
}

func calculateLoadInColumn(column string) int {
	load := 0

	for i, v := range column {
		if v == 'O' {
			load += len(column) - i
		}
	}

	return load
}

func hash(rows []string) string {
	h := sha256.New()
	for _, row := range rows {
		h.Write([]byte(row))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func spinCycle(rows []string) {
	var freeIndex int
	var newRow string
	for colIndex := range rows[0] {
		// tilt this column North
		freeIndex = 0

		for rowIndex := range rows {
			switch rows[rowIndex][colIndex] {
			case 'O':
				rows[freeIndex] = rows[freeIndex][:colIndex] + "O" + rows[freeIndex][colIndex+1:]
				freeIndex++
			case '#':
				for i := freeIndex; i < rowIndex; i++ {
					rows[i] = rows[i][:colIndex] + "." + rows[i][colIndex+1:]
				}
				freeIndex = rowIndex + 1
			}
		}

		for i := freeIndex; i < len(rows); i++ {
			rows[i] = rows[i][:colIndex] + "." + rows[i][colIndex+1:]
		}
	}

	for rowIndex := range rows {
		// tilt this row West
		freeIndex = 0
		newRow = ""

		for colIndex, v := range rows[rowIndex] {
			switch v {
			case 'O':
				newRow += "O"
				freeIndex++
			case '#':
				newRow += strings.Repeat(".", colIndex-freeIndex)
				newRow += "#"
				freeIndex = colIndex + 1
			}
		}
		newRow += strings.Repeat(".", len(rows[0])-freeIndex)
		rows[rowIndex] = newRow
	}

	for colIndex := range rows[0] {
		// tilt this column South
		freeIndex = len(rows) - 1

		for rowIndex := len(rows) - 1; rowIndex >= 0; rowIndex-- {
			switch rows[rowIndex][colIndex] {
			case 'O':
				rows[freeIndex] = rows[freeIndex][:colIndex] + "O" + rows[freeIndex][colIndex+1:]
				freeIndex--
			case '#':
				for i := freeIndex; i > rowIndex; i-- {
					rows[i] = rows[i][:colIndex] + "." + rows[i][colIndex+1:]
				}
				freeIndex = rowIndex - 1
			}

		}

		for i := freeIndex; i >= 0; i-- {
			rows[i] = rows[i][:colIndex] + "." + rows[i][colIndex+1:]
		}
	}

	for rowIndex := range rows {
		// tilt this row East
		freeIndex = len(rows[0]) - 1
		newRow = ""

		for colIndex := freeIndex; colIndex >= 0; colIndex-- {
			switch rows[rowIndex][colIndex] {
			case 'O':
				newRow = "O" + newRow
				freeIndex--
			case '#':
				newRow = strings.Repeat(".", freeIndex-colIndex) + newRow
				newRow = "#" + newRow
				freeIndex = colIndex - 1
			}
		}

		newRow = strings.Repeat(".", freeIndex+1) + newRow
		rows[rowIndex] = newRow
	}
}

func displayRows(rows []string) {
	fmt.Println("======================")
	for _, v := range rows {
		fmt.Println(v)
	}
}

func main() {
	lines := fetchLines(14)

	////////////////////////////////////////////////

	sum := 0
	var column string
	for i := range lines[0] {
		column = ""
		for j := range lines {
			column += string(lines[j][i])
		}
		sum += tiltNorthAndCalculateLoadInColumn(column)
	}

	fmt.Println(sum)

	////////////////////////////////////////////////

	// hash of the curent state to the number of spin cycles it took us to get there
	hashToNumSpinCycles := make(map[string]int)

	// perform 1000000000 spin cycles
	fastForward := false
	spinCyclesPerformed := 0
	for spinCyclesPerformed != 1000000000 {
		currentConfigHash := hash(lines)
		numSpinCycles, exists := hashToNumSpinCycles[currentConfigHash]
		if exists && !fastForward {
			// this means that currentConfigHash and numSpinCycles are the number of spin cycles
			// which produce the SAME hash
			// so we can sefely add on the difference of these
			difference := spinCyclesPerformed - numSpinCycles

			// so it turns out that as long as our num of spin cycles modulo the difference is the same, we will be in the same state

			// see how close we can get to the target number of spin cycles and then fast forward there
			bigSpinCycles := 1000000000
			for bigSpinCycles%difference != spinCyclesPerformed%difference {
				bigSpinCycles--
			}
			spinCyclesPerformed = bigSpinCycles
			// only fast forward once
			// from now on, we will have to perform spin cycles manually
			fastForward = true
			continue
		}
		hashToNumSpinCycles[currentConfigHash] = spinCyclesPerformed

		spinCycle(lines)
		spinCyclesPerformed++
	}

	// calculate load in columns
	sum = 0
	for i := range lines[0] {
		column = ""
		for j := range lines {
			column += string(lines[j][i])
		}
		sum += calculateLoadInColumn(column)
	}
	fmt.Println(sum)
}
