build:
	go build ./cmd/query-counter
run:
	./query-counter -n 5 -ifile "./files/input.txt" -ofile "./files/output.tsv" -type approximate
run-old:
	./query-counter -n 5 -ifile "./files/input.txt" -ofile "./files/exact_output.tsv" -type precision
