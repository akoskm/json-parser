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
		fmt.Println(error)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	insideName := false
	insideValue := false
	name := ""
	value := ""
	previousChar := ""
	invalid := false
	lineCount := 1
	var objectStack []bool
	var arrayStack []bool
	for {
		if invalid {
			msg := fmt.Sprintf("unexpected character %v on line %v", previousChar, lineCount)
			return false, errors.New(msg)
		}
		readRune, _, err := reader.ReadRune()

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			} else {
				fmt.Println("finished parsing")
			}
			if len(objectStack) != 0 {
				return false, errors.New("object is not closed")
			}
			if len(arrayStack) != 0 {
				return false, errors.New("array is not closed")
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
			objectStack = append(objectStack, true)
		case "}":
			if previousChar == "," {
				fmt.Println("no, after the last object")
				invalid = true
			}
			if len(objectStack) > 0 {
				objectStack = objectStack[:len(objectStack)-1]
			} else {
				fmt.Println("no matching {")
				invalid = true
			}
		case "[":
			arrayStack = append(arrayStack, true)
		case "]":
			if previousChar == "," {
				fmt.Println("no, after the last object")
				invalid = true
			}
			if len(arrayStack) > 0 {
				arrayStack = arrayStack[:len(arrayStack)-1]
			} else {
				fmt.Println("no matching {")
				invalid = true
			}
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
			fmt.Println("found name or value", result)
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
						fmt.Println("got a boolean", guessBool)
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
						fmt.Println("got a null", guessNull)
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
						fmt.Printf("%q looks like a number.\n", guessNumber)
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
	return true, nil
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("must provide a filename")
		os.Exit(1)
	}

	filename := args[len(args)-1]

	_, err := parseJSON2(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
