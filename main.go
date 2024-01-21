package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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
	for {
		readRune, _, err := reader.ReadRune()

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		char := string(readRune)

		if char == "{" {
			beginObject = true
		}
		if char == "}" {
			endObject = true
		}
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

	parseJSON(filename)
}
