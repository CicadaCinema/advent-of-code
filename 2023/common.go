package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
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
