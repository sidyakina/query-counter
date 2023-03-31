package app

import (
	"fmt"
	"strconv"
	"strings"
)

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
