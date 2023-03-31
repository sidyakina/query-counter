package main

import (
	"flag"
	"log"
)

func main() {
	var n uint
	var inputFile, outputFile string

	flag.UintVar(&n, "n", 0, "N - max queries in memory")
	flag.StringVar(&inputFile, "ifile", "input.txt", "input file, default: input.txt")
	flag.StringVar(&outputFile, "ofile", "output.tsv", "output file, default: output.tsv")

	flag.Parse()

	if n == 0 {
		log.Panic("max queries value is required")
	}

	log.Printf("%v, %v, %v", n, inputFile, outputFile)
}
