package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

func parseLines(records [][]string) []problem {
	problems := make([]problem, len(records))

	for i, val := range records {
		problems[i] = problem{
			question: strings.TrimSpace(val[0]),
			answer:   strings.TrimSpace(val[1]),
		}
	}

	return problems
}

func readFile(r io.Reader) ([][]string, error) {
	// read the file line by line
	var records [][]string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		records = append(records, strings.Split(line, ","))
	}

	return records, scanner.Err()
}

func getUserAnswer(i int, p problem) (isCorrect bool) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Problem #%d: %s = ", i+1, p.question)
	scanner.Scan()

	inputAnswer := scanner.Text()
	inputAnswer = strings.TrimSpace(inputAnswer)

	return inputAnswer == p.answer
}

func main() {
	// parse flags for the CSV filename
	csvFile := flag.String("csv", "problems.csv", "a csv file in the form of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatalf("Failed to open the CSV file %q", *csvFile)
	}
	defer file.Close()

	records, err := readFile(file)
	if err != nil {
		log.Fatalf("Failed to read the CSV file %q", *csvFile)
	}

	problems := parseLines(records)

	var score int
	for i, p := range problems {
		if isCorrect := getUserAnswer(i, p); isCorrect {
			score++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", score, len(problems))
}
