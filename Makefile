build:
	go build ./cmd/query-counter
run:
	./query-counter -n 50 -ifile "./files/input.txt" -ofile "./files/output.tsv"