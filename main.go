package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

func parseJSON2(filename string) (bool, error) {
	file, error := os.Open(filename)
	if error != nil {
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	beginObject := false
	endObject := false
	insideName := false
	insideValue := false
	name := ""
	value := ""
	previousChar := ""
	invalid := false
	lineCount := 1
	for {
		if invalid {
			msg := fmt.Sprintf("unexpected character %v on line %v", previousChar, lineCount)
			return false, errors.New(msg)
		}
		readRune, _, err := reader.ReadRune()

		if err != nil {
			if err != io.EOF {
			} else {
			}
			break
		}

		char := string(readRune)

		if char == "\n" {
			lineCount++
		}

		if insideName {
			name += char
		}
		if insideValue {
			value += char
		}

		switch char {
		case "{":
			beginObject = true
		case "}":
			if previousChar == "," {
				invalid = true
			}
			endObject = true
		case "\"":
			result := ""
			for {
				nameOrValue, _, _ := reader.ReadRune()
				parsedChar := string(nameOrValue)
				if parsedChar == "\"" {
					break
				}
				result += parsedChar
			}
		default:
			if unicode.IsSpace(readRune) {
				continue
			}
			if unicode.IsDigit(readRune) || unicode.IsLetter(readRune) {
				if char == "t" || char == "f" {
					guessBool := char
					for {
						gbCharRaw, _, _ := reader.ReadRune()
						gbChar := string(gbCharRaw)
						if !unicode.IsLetter(gbCharRaw) {
							break
						}
						guessBool += gbChar
					}
					if guessBool == "true" || guessBool == "false" {
						continue
					}
				} else if char == "n" {
					guessNull := char
					for {
						gbCharRaw, _, _ := reader.ReadRune()
						gbChar := string(gbCharRaw)
						if !unicode.IsLetter(gbCharRaw) {
							break
						}
						guessNull += gbChar
					}
					if guessNull == "null" {
						continue
					}
				} else if unicode.IsDigit(readRune) {
					guessNumber := char
					for {
						gbCharRaw, _, _ := reader.ReadRune()
						gbChar := string(gbCharRaw)
						if !unicode.IsDigit(gbCharRaw) {
							break
						}
						guessNumber += gbChar
					}
					// check if guessNumber is a number
					if _, err := strconv.Atoi(guessNumber); err == nil {
						continue
					}
				}
				// letters and symbols only allowed inside ""
				if !insideName && !insideValue {
					invalid = true
				}
			}
		}
		previousChar = char
	}
	if beginObject == false || endObject == false {
	}
	return true, nil
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		os.Exit(1)
	}

	filename := args[len(args)-1]

	_, err := parseJSON2(filename)
	if err != nil {
		os.Exit(1)
	}
}
