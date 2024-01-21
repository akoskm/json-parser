package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

func parseJSON2(filename string) {
	file, error := os.Open(filename)
	if error != nil {
		fmt.Println(error)
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
			fmt.Println(fmt.Sprintf("unexpected character %v on line %v", previousChar, lineCount))
			os.Exit(1)
		}
		readRune, _, err := reader.ReadRune()

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			} else {
				fmt.Println("finished parsing")
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
				fmt.Println("no, after the last object")
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
			fmt.Println("found name or value", result)
		case ",":
			if previousChar != "\"" {
				invalid = true
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
		fmt.Println("invalid JSON")
	}

}

func parseJSON(filename string) {
	file, error := os.Open(filename)
	if error != nil {
		fmt.Println(error)
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
	valueComes := false
	previousChar := ""
	invalid := false
	foundOpeningQuote := false
	lineCount := 1
	for {
		if invalid {
			fmt.Println(fmt.Sprintf("unexpected character %v on line %v", previousChar, lineCount))
			os.Exit(1)
		}
		readRune, _, err := reader.ReadRune()

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			} else {
				fmt.Println("finished parsing")
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

		fmt.Println("char", char)
		switch char {
		case "{":
			beginObject = true
		case "}":
			if previousChar == "," {
				fmt.Println("no ,  after the last object")
				invalid = true
			}
			endObject = true
		case "\"":
			if !foundOpeningQuote {
				if valueComes && !insideValue {
					insideValue = true
				} else if !insideValue {
					insideName = true
				}
			} else {
				if insideName {
					insideName = false
				} else if insideValue {
					insideValue = false
					valueComes = false
				}
				foundOpeningQuote = false
			}
		case ":":
			valueComes = true
		case ",":
			if insideValue || insideName {
				return
			}
			if previousChar != "\"" {
				invalid = true
			}
		default:
			if unicode.IsSpace(readRune) {
				continue
			}
			fmt.Println("default case", insideName, insideValue)
			if unicode.IsDigit(readRune) || unicode.IsLetter(readRune) {
				// letters and symbols only allowed inside ""
				if !insideName && !insideValue {
					invalid = true
				}
			}
		}
		previousChar = char
	}
	if beginObject == false || endObject == false {
		fmt.Println("invalid JSON")
	}
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("must provide a filename")
		os.Exit(1)
	}

	filename := args[len(args)-1]

	parseJSON2(filename)
}
