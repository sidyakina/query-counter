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
	for scanner.Scan() {
		query := scanner.Text()

		exists := counter.table.Add(query)
		if !exists {
			_, err = counter.dictionaryFile.WriteString(query + "\n")
			if err != nil {
				return fmt.Errorf("failed to write query to dictionary file: %w", err)
			}
		}
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
	table          *Table
	inputFile      *os.File
	dictionaryFile *os.File
	outputFile     *os.File
}

func (c *Counter) setup(n int, inputFilePath string) error {
	c.table = NewTable(n)

	var err error
	// source file
	c.inputFile, err = os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}

	// temp file for discarded queries
	c.dictionaryFile, err = os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
	_, err = c.dictionaryFile.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to rewind dictionary file: %w", err)
	}

	return nil
}

func (c *Counter) outputResult() error {
	scanner := bufio.NewScanner(c.dictionaryFile)
	for scanner.Scan() {
		query := scanner.Text()
		count := c.table.Count(query)

		_, err := c.outputFile.WriteString(formatQueryData(query, count) + "\n")
		if err != nil {
			return fmt.Errorf("failed to write query to output file: %w", err)
		}
	}

	return nil
}

func formatQueryData(query string, count int32) string {
	return fmt.Sprintf("%v\t%v", query, count)
}

func clean(counter *Counter) {
	for _, file := range []*os.File{counter.inputFile, counter.dictionaryFile, counter.outputFile} {
		if file != nil {
			err := file.Close()
			if err != nil {
				log.Printf("failed to close file %v: %v", file.Name(), err)
			}
		}
	}

	err := os.Remove(tempFile)
	if err != nil {
		log.Printf("failed to remove temp file: %v", err)
	}
}
