package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	input := fetchInput(2)

	sum := 0
	power := 0
	rGameId := regexp.MustCompile("[0-9]+")
	rRed := regexp.MustCompile("([0-9]+) red")
	rGreen := regexp.MustCompile("([0-9]+) green")
	rBlue := regexp.MustCompile("([0-9]+) blue")
	var r, g, b int
	var minR, minG, minB int
	var fail bool
	var err error

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		fail = false
		minR, minG, minB = 0, 0, 0

		for _, set := range strings.Split(line, ";") {
			r, g, b = 0, 0, 0

			if match := rRed.FindStringSubmatch(set); len(match) == 2 {
				r, err = strconv.Atoi(match[1])
				check(err)
			} else if len(match) != 0 {
				panic("invalid number of matches")
			}
			if match := rGreen.FindStringSubmatch(set); len(match) == 2 {
				g, err = strconv.Atoi(match[1])
				check(err)
			} else if len(match) != 0 {
				panic("invalid number of matches")
			}
			if match := rBlue.FindStringSubmatch(set); len(match) == 2 {
				b, err = strconv.Atoi(match[1])
				check(err)
			} else if len(match) != 0 {
				panic("invalid number of matches")
			}

			if r > minR {
				minR = r
			}
			if g > minG {
				minG = g
			}
			if b > minB {
				minB = b
			}

			if r > 12 || g > 13 || b > 14 {
				fail = true
			}
		}

		power += minR * minG * minB

		if fail {
			continue
		}

		id, err := strconv.Atoi(rGameId.FindString(line))
		check(err)
		sum += id
	}

	fmt.Println(sum)
	fmt.Println(power)
}
