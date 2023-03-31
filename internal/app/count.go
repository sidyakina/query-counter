package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const tempDir = "./temp"

func Count(n uint, inputFilePath, outputFilePath string) error {
	// we use n-1 because query will be read before add in chan (so we have n-1 in chan and one is waiting)
	queries := make(chan string, n-1)
	var readErr error

	go func() {
		// we can read our queries while handle previous ones, it isn't necessary for this case
		// but can make our Count faster if queries' source is another
		readErr = readQueries(inputFilePath, queries)
		close(queries)
	}()

	err := os.MkdirAll(tempDir, 0777)
	if err != nil {
		return fmt.Errorf("create dir temp: %w", err)
	}
	defer clear()

	for query := range queries {
		err := handleQuery(query)
		if err != nil {
			return fmt.Errorf("failed to handle query: %w", err)
		}
	}

	if readErr != nil {
		return fmt.Errorf("failed to read data from %v: %w", inputFilePath, readErr)
	}

	err = makeOutputFile(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to make output file: %w", err)
	}

	return nil
}

func clear() {
	err := os.RemoveAll(tempDir)
	if err != nil {
		log.Printf("failed to remove dir temp: %v", err)
	}
}

func readQueries(inputFilePath string, readCh chan string) error {
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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		q := scanner.Text()
		readCh <- q
	}

	return nil
}

func makeOutputFile(fileName string) error {
	outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %w", fileName, err)
	}

	tempFiles, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read dir %v: %w", tempDir, err)
	}

	for _, tempFile := range tempFiles {
		data, err := os.ReadFile(fmt.Sprintf("%v/%v", tempDir, tempFile.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file %v: %w", tempFile.Name(), err)
		}

		_, err = outputFile.WriteString(string(data) + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file %v: %w", fileName, err)
		}
	}

	return nil
}
