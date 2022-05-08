package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
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

func checkUserAnswer(i int, p problem) (isCorrect bool) {
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
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
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

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	defer timer.Stop()

	score := 0

questionLoop:
	for i, p := range problems {
		resultChan := make(chan bool)

		go func(i int, p problem) { // goroutine to check the user's answer
			resultChan <- checkUserAnswer(i, p)
		}(i, p)

		select {
		case result := <-resultChan:
			if result {
				score++
			}
		case <-timer.C:
			fmt.Println("\nYou ran out of time!")
			break questionLoop
		}
	}

	fmt.Printf("You scored %d out of %d.\n", score, len(problems))
}
