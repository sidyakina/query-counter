package precision

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"os"
)

// for each query we create file for count using hash
// if two queries have same hash "i" used for avoid collision
func handleQuery(query string) error {
	h := hash(query)

	var fileName string
	var count int64

	for i := 0; ; i++ {
		fileName = fmt.Sprintf("%v/%v(%v).txt", tempDir, h, i)
		data, err := os.ReadFile(fileName)

		if errors.Is(err, os.ErrNotExist) {
			// it's first occurrence of query
			count = 1
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read file %v: %w", fileName, err)
		}

		queryInFile, prevNumber, err := parseQueryData(string(data))
		if err != nil {
			return fmt.Errorf("failed to parse file %v: %w", fileName, err)
		}

		if queryInFile == query {
			// it isn't first occurrence of query
			count = prevNumber + 1

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

	_, err = file.Write([]byte(formatQueryData(query, count)))
	if err != nil {
		return fmt.Errorf("failed to write to file %v: %w", fileName, err)
	}

	return nil
}

func hash(query string) string {
	return fmt.Sprintf("%x", md5.New().Sum([]byte(query)))
}
