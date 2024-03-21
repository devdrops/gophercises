package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	// CSV file
	csvFile := flag.String("csv", "problems.csv", "CSV file in format 'question,answer'.")
	// Timeout
	timeout := flag.Int("timeout", 30, "Timeout limit to give an answer, in seconds.")
	flag.Parse()

	// Reading file
	file, err := os.Open(*csvFile)
	if err != nil {
		fmt.Printf("Failed to open %s (%s)\n", *csvFile, err.Error())
		os.Exit(1)
	}
	// Close after everything
	defer file.Close()

	var totals, correct, wrong int

	timer := time.NewTimer(time.Duration(*timeout) * time.Second)

	// Iterate data
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Parse CSV line
		r := csv.NewReader(strings.NewReader(scanner.Text()))
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		// Display question
		fmt.Printf("Q: %s = ? :\t", line[0])
		ac := make(chan string)
		go func() {
			// Read user input
			var input string
			fmt.Scanf("%s\n", &input)
			ac <- input
		}()

		select {
		case <-timer.C:
			// Final score
			fmt.Println()
			fmt.Println(strings.Repeat("-", 60))
			fmt.Printf("SCORE: %d Questions\nRight answers: %d, Wrong answers: %d\n", totals, correct, wrong)
			os.Exit(0)
		case answer := <- ac:
			// Results track
			if answer == line[1] {
				correct++
			} else {
				wrong++
			}

			// Display result
			fmt.Printf("A: %s, %t\n", answer, (answer == line[1]))

			totals++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
