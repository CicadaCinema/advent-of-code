package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	input := fetchInput(1)

	////////////////////////////////////////////////

	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		chars := strings.Split(line, "")

		// search for the first digit
		for _, char := range chars {
			firstDigit, err := strconv.Atoi(char)
			if err == nil {
				sum += 10 * firstDigit
				break
			}
		}

		// search for the last digit
		for i := len(chars) - 1; i >= 0; i-- {
			lastDigit, err := strconv.Atoi(chars[i])
			if err == nil {
				sum += lastDigit
				break
			}
		}
	}
	fmt.Println(sum)

	////////////////////////////////////////////////

	sum = 0
	r1 := regexp.MustCompile("one|two|three|four|five|six|seven|eight|nine|[1-9]")
	r2 := regexp.MustCompile("(.*)(one|two|three|four|five|six|seven|eight|nine|[1-9])")
	for _, line := range lines {
		// the input ends with an empty line
		if line == "" {
			continue
		}

		// search for the first digit
		switch r1.FindString(line) {
		case "1", "one":
			sum += 10
		case "2", "two":
			sum += 20
		case "3", "three":
			sum += 30
		case "4", "four":
			sum += 40
		case "5", "five":
			sum += 50
		case "6", "six":
			sum += 60
		case "7", "seven":
			sum += 70
		case "8", "eight":
			sum += 80
		case "9", "nine":
			sum += 90
		}

		// search for the last digit
		switch r2.FindStringSubmatch(line)[2] {
		case "1", "one":
			sum += 1
		case "2", "two":
			sum += 2
		case "3", "three":
			sum += 3
		case "4", "four":
			sum += 4
		case "5", "five":
			sum += 5
		case "6", "six":
			sum += 6
		case "7", "seven":
			sum += 7
		case "8", "eight":
			sum += 8
		case "9", "nine":
			sum += 9
		}
	}
	fmt.Println(sum)
}
