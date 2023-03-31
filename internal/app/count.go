package app

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

func handleQuery(query string) error {
	h := hash(query)

	var fileName string
	var number int64

	for i := 0; ; i++ {
		fileName = fmt.Sprintf("%v/%v(%v).txt", tempDir, h, i)
		data, err := os.ReadFile(fileName)

		if errors.Is(err, os.ErrNotExist) {
			// it's first occurrence of query
			number = 1
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read file %v: %w", fileName, err)
		}

		queryInFile, prevNumber, err := parseDataInFile(string(data))
		if err != nil {
			return fmt.Errorf("failed to parse file %v: %w", fileName, err)
		}

		if queryInFile == query {
			// it isn't first occurrence of query
			number = prevNumber + 1

			break
		}

		// it's possible if hash(query1) == hash(query2), we need to check other files or create file for this query
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %w", fileName, err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close input file: %v", err)
		}
	}()

	_, err = file.Write([]byte(formatDataForQuery(query, number)))
	if err != nil {
		return fmt.Errorf("failed to write to file %v: %w", fileName, err)
	}

	return nil
}

func formatDataForQuery(query string, number int64) string {
	return fmt.Sprintf("%v\t%v", query, number)
}

func parseDataInFile(data string) (string, int64, error) {
	// data format: "{query}/t{number}"

	temp := strings.Split(data, "\t")
	if len(temp) != 2 {
		return "", 0, fmt.Errorf("wrong data in file: %v", temp)
	}

	number, err := strconv.ParseInt(temp[1], 10, 64)
	if err != nil {
		return "", 0, err
	}

	return temp[0], number, nil
}

func hash(query string) string {
	return string(query[0])
}
