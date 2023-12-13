package main

import (
	"fmt"
)

func solveSmudge(block []string) int {
	prevValue := solve(block, 0)

	// we want to identify a smudge which does not return the prevValue
	var originalRow string
	var trialSol int
	for i := range block {
		originalRow = block[i]
		for j := range block[i] {
			switch block[i][j] {
			case '#':
				block[i] = block[i][:j] + "." + block[i][j+1:]
			case '.':
				block[i] = block[i][:j] + "#" + block[i][j+1:]
			default:
				panic("unexpected rune")
			}

			trialSol = solve(block, prevValue)
			if trialSol != 0 {
				return trialSol
			}

			// return this row to its original state before we try the next smudge
			block[i] = originalRow
		}
	}

	panic("no smudge caused us to return a different value")
}

// skips the solution with value skipSolution
// if no solution can be found, returns 0
func solve(block []string, skipSolution int) int {
	// each element of this slice is a summary of a single column
	columnSummaries := make([]string, len(block[0]))
	// the row summaries are block itself!!

	var columnSummary string
	for i := range columnSummaries {
		columnSummary = ""
		for j := range block {
			columnSummary += string(block[j][i])
		}
		columnSummaries[i] = columnSummary
	}

	var leftSide []string
	var lenLeftSide int
	var rightSide []string
	var lenRightSide int
	var minLength int
outerC:
	for columnIndex := 1; columnIndex <= len(columnSummaries)-1; columnIndex++ {
		// see if there is a reflection between this index and the next one along
		leftSide = columnSummaries[:columnIndex]
		lenLeftSide = len(leftSide)
		rightSide = columnSummaries[columnIndex:]
		lenRightSide = len(rightSide)

		minLength = lenLeftSide
		if lenRightSide < minLength {
			minLength = lenRightSide
		}

		l := lenLeftSide - 1
		r := 0
		for i := 0; i < minLength; i++ {
			if leftSide[l] != rightSide[r] {
				continue outerC
			}
			l--
			r++
		}

		// there is a reflection, so return the 1-indexed result
		if columnIndex != skipSolution {
			return columnIndex
		}
	}

outerR:
	for rowIndex := 1; rowIndex <= len(block)-1; rowIndex++ {
		leftSide = block[:rowIndex]
		lenLeftSide = len(leftSide)
		rightSide = block[rowIndex:]
		lenRightSide = len(rightSide)

		minLength = lenLeftSide
		if lenRightSide < minLength {
			minLength = lenRightSide
		}

		l := lenLeftSide - 1
		r := 0
		for i := 0; i < minLength; i++ {
			if leftSide[l] != rightSide[r] {
				continue outerR
			}
			l--
			r++
		}

		if rowIndex*100 != skipSolution {
			return rowIndex * 100
		}
	}

	return 0
}

func main() {
	lines := fetchLines(13)

	////////////////////////////////////////////////

	sum := 0
	start := 0
	for i, line := range lines {
		if line == "" {
			sum += solve(lines[start:i], 0)
			start = i + 1
		}
	}
	sum += solve(lines[start:], 0)

	fmt.Println(sum)

	////////////////////////////////////////////////

	sum = 0
	start = 0
	for i, line := range lines {
		if line == "" {
			sum += solveSmudge(lines[start:i])
			start = i + 1
		}
	}
	sum += solveSmudge(lines[start:])

	fmt.Println(sum)
}
