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
	counter := &Counter{}
	defer clean(counter)

	err := counter.setup(n, inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to setup counter: %v", err)
	}

	scanner := bufio.NewScanner(counter.inputFile)
	index := 0
	for scanner.Scan() {
		q := scanner.Text()

		sampleNumber := index / n
		discardedQuery, count := counter.queries.Add(q, sampleNumber)
		if discardedQuery != "" {
			_, err = counter.discardFile.WriteString(formatQueryData(discardedQuery, count) + "\n")
			if err != nil {
				return fmt.Errorf("failed to write discarded query to temp file: %w", err)
			}
		}

		index += 1
	}

	err = counter.setupOutput(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to setup counter before output: %v", err)
	}

	err = counter.outputResult()
	if err != nil {
		return fmt.Errorf("failed to output results: %v", err)
	}

	return nil
}

type Counter struct {
	queries     *Queries
	inputFile   *os.File
	discardFile *os.File
	outputFile  *os.File
}

func (c *Counter) setup(n int, inputFilePath string) error {
	c.queries = NewQueries(int(n) - 1) // because we will have capacity + 1 while adding

	var err error
	// source file
	c.inputFile, err = os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}

	// temp file for discarded queries
	c.discardFile, err = os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to create temp file %v: %w", tempFile, err)
	}

	return nil
}

func (c *Counter) setupOutput(outputFilePath string) error {
	var err error
	c.outputFile, err = os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to open output file %v: %w", outputFilePath, err)
	}

	// rewind file for reading
	_, err = c.discardFile.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to rewind discard file: %w", err)
	}

	return nil
}

func (c *Counter) outputResult() error {
	var err error

	for k, v := range c.queries.queries {
		_, err = c.outputFile.WriteString(formatQueryData(k, v.count) + "\n")
		if err != nil {
			return fmt.Errorf("failed to write query to output file: %w", err)
		}
	}

	scanner := bufio.NewScanner(c.discardFile)
	for scanner.Scan() {
		_, err = c.outputFile.WriteString(scanner.Text() + "\n")
		if err != nil {
			return fmt.Errorf("failed to write query to output file: %w", err)
		}
	}

	return nil
}

func formatQueryData(query string, number int) string {
	return fmt.Sprintf("%v\t%v", query, number)
}

func clean(counter *Counter) {
	for _, file := range []*os.File{counter.inputFile, counter.discardFile, counter.outputFile} {
		if file != nil {
			err := file.Close()
			if err != nil {
				log.Printf("failed to close file %v: %v", file.Name(), err)
			}
		}
	}

	err := os.RemoveAll(tempFile)
	if err != nil {
		log.Printf("failed to remove temp file: %v", err)
	}
}
