package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

func splitWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func countOccurrencesSequential(records [][]string, targetWord string) int {
	totalCount := 0

	for _, row := range records {
		rowText := strings.Join(row, " ")
		totalCount += countOccurrences(rowText, targetWord)
	}

	return totalCount
}

func countOccurrencesConcurrent(records [][]string, targetWord string) int {
	totalCount := 0
	ch := make(chan int, len(records))

	for _, row := range records {
		go func(rowText string) {
			count := countOccurrences(rowText, targetWord)
			ch <- count
		}(strings.Join(row, " "))
	}

	for range records {
		count := <-ch
		totalCount += count
	}

	close(ch)

	return totalCount
}

func countOccurrences(text string, targetWord string) int {
	text = strings.ToLower(text)
	targetWord = strings.ToLower(targetWord)

	words := splitWords(text)

	count := 0

	for _, word := range words {
		if strings.Contains(word, targetWord) {
			count++
		}
	}

	return count
}

func main() {
	wordPtr := flag.String("word", "", "the word to count")
	concurrentPtr := flag.Bool("concurrent", false, "perform concurrent word count")
	flag.Parse()

	if *wordPtr == "" {
		fmt.Println("Please provide the word to count using the -word flag.")
		return
	}

	// Open the CSV file
	file, err := os.Open("testcur.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	var totalCount int
	var duration time.Duration

	if *concurrentPtr {
		// Concurrent word count
		startTimeConcurrent := time.Now()
		totalCount = countOccurrencesConcurrent(records, *wordPtr)
		duration = time.Since(startTimeConcurrent)
		fmt.Printf("Concurrent word count: The word '%s' occurs %d times.\n", *wordPtr, totalCount)
		fmt.Printf("Time taken for concurrent word count: %s\n", duration)
	} else {
		// Sequential word count
		startTimeSequential := time.Now()
		totalCount = countOccurrencesSequential(records, *wordPtr)
		duration = time.Since(startTimeSequential)
		fmt.Printf("Sequential word count: The word '%s' occurs %d times.\n", *wordPtr, totalCount)
		fmt.Printf("Time taken for sequential word count: %s\n", duration)
	}
}
