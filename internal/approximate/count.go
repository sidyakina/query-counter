package approximate

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const tempFile = "temp.txt"

func Count(n int, inputFilePath, outputFilePath string) error {
	queries := NewQueries(int(n) - 1) // because we will have capacity + 1 while adding

	// source file
	file, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close input file: %v", err)
		}
	}()

	// temp file for discarded queries
	discardFile, err := os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to create temp file %v: %w", tempFile, err)
	}

	defer func() {
		err := discardFile.Close()
		if err != nil {
			log.Printf("failed to close input file: %v", err)
		}

		clean()
	}()

	scanner := bufio.NewScanner(file)
	index := 0
	for scanner.Scan() {
		q := scanner.Text()

		sampleNumber := index / n
		discardedQuery, count := queries.Add(q, sampleNumber)
		if discardedQuery != "" {
			_, err = discardFile.WriteString(formatQueryData(discardedQuery, count) + "\n")
			if err != nil {
				return fmt.Errorf("failed to write discarded query to temp file: %w", err)
			}
		}

		index += 1
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %w", outputFilePath, err)
	}

	defer func() {
		err := outputFile.Close()
		if err != nil {
			log.Printf("failed to close input file: %v", err)
		}
	}()

	for k, v := range queries.queries {
		_, err = outputFile.WriteString(formatQueryData(k, v.count) + "\n")
		if err != nil {
			return fmt.Errorf("failed to write query to output file: %w", err)
		}
	}

	_, err = discardFile.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to rewind discard file: %w", err)
	}

	scanner = bufio.NewScanner(discardFile)
	for scanner.Scan() {
		_, err = outputFile.WriteString(scanner.Text() + "\n")
		if err != nil {
			return fmt.Errorf("failed to write query to output file: %w", err)
		}
	}

	return nil
}

func clean() {
	err := os.RemoveAll(tempFile)
	if err != nil {
		log.Printf("failed to remove temp file: %v", err)
	}
}

func formatQueryData(query string, number int) string {
	return fmt.Sprintf("%v\t%v", query, number)
}
