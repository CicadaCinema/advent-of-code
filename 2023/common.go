package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// https://github.com/CicadaCinema/six-small-algorithms-go/blob/d705e655f18044650dd906dab50a1703bf469d96/5/gcd_lcm.go#L8-L45
func gcd(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	if a == b {
		return a
	}

	if a < 0 {
		a *= -1
	}
	if b < 0 {
		b *= -1
	}

	var A, B int
	if a > b {
		A = a
		B = b
	} else {
		A = b
		B = a
	}

	R := A % B

	return gcd(B, R)
}

func lcm(a, b int) int {
	return (a / gcd(a, b)) * b
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return +1
	}
	return 0
}

func assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func uintegers(line string) []int {
	rNum := regexp.MustCompile("[0-9]+")
	stringResult := rNum.FindAllString(line, -1)
	result := make([]int, len(stringResult))
	for i := range stringResult {
		v, err := strconv.Atoi(stringResult[i])
		check(err)
		result[i] = v
	}
	return result
}

func integers(line string) []int {
	rNum := regexp.MustCompile("-?[0-9]+")
	stringResult := rNum.FindAllString(line, -1)
	result := make([]int, len(stringResult))
	for i := range stringResult {
		v, err := strconv.Atoi(stringResult[i])
		check(err)
		result[i] = v
	}
	return result
}

func sliceSum(s []int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func fetchLines(day int) []string {
	input := fetchInput(day)
	lines := strings.Split(input, "\n")

	if lines[len(lines)-1] != "" {
		panic("last line of input data was not empty as expected")
	}

	return lines[0 : len(lines)-1]
}

func fetchInput(day int) string {
	// validate input
	if day < 1 || day > 25 {
		panic("invalid day")
	}

	localInputFilePath := fmt.Sprintf("input-%02d.txt", day)

	if _, err := os.Stat(localInputFilePath); err == nil {
		// cached input file exists
		cachedInput, err := os.ReadFile(localInputFilePath)
		check(err)
		return string(cachedInput)
	} else if !errors.Is(err, os.ErrNotExist) {
		// there was an error different to 'file does not exist'
		panic(err)
	}

	session, err := os.ReadFile("session.txt")
	check(err)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/2023/day/%d/input", day), nil)
	check(err)

	req.AddCookie(&http.Cookie{Name: "session", Value: string(session)})

	resp, err := http.DefaultClient.Do(req)
	check(err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	check(err)

	inputString := string(body)

	// save the cached input
	err = os.WriteFile(localInputFilePath, body, 0644)
	check(err)

	return inputString
}
