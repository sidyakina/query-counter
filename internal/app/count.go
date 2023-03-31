package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const tempDir = "./temp"

func Count(n uint, inputFilePath, outputFilePath string) error {
	_, _ = n, outputFilePath

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

	err = os.MkdirAll(tempDir, 0777)
	if err != nil {
		return fmt.Errorf("create dir temp: %w", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		q := scanner.Text()
		log.Println(q)
		err := handleQuery(q)
		if err != nil {
			return fmt.Errorf("failed to handle query: %w", err)
		}
	}

	err = makeOutputFile(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to make output file: %w", err)
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		// output file was created successfully, print as warning
		log.Printf("failed to remove dir temp: %v", err)
	}

	return nil
}

func makeOutputFile(fileName string) error {
	outputFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %w", fileName, err)
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read dir %v: %w", tempDir, err)
	}

	for _, tempFile := range files {
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
