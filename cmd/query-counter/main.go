package main

import (
	"flag"
	"fmt"
	"log"
	"query-counter/internal/approximate"
	"query-counter/internal/precision"
)

func main() {
	var n uint
	var inputFile, outputFile, counterType string

	flag.UintVar(&n, "n", 0, "N - max queries in memory")
	flag.StringVar(&inputFile, "ifile", "input.txt", "input file, default: input.txt")
	flag.StringVar(&outputFile, "ofile", "output.tsv", "output file, default: output.tsv")
	flag.StringVar(&counterType, "type", "approximate", "counter type, default: approximate")

	flag.Parse()

	if n == 0 {
		log.Panic("max queries value is required")
	}

	log.Printf("count with counter type: %v", counterType)

	var err error
	switch counterType {
	case "precision":
		// old slow counter
		err = precision.Count(n, inputFile, outputFile)
	case "approximate":
		err = approximate.Count(int(n), inputFile, outputFile)
	default:
		err = fmt.Errorf("unknown counter type: %v", counterType)
	}

	if err != nil {
		log.Panicf("failed to count: %v", err)
	}
}
