# Query-Counter

## Task
Нужно написать консольную golang утилиту, не используя сторонние библиотеки и базы
данных. На входе файл с поисковыми запросами, на выходе файл в формате tsv с
уникальными запросами и их приблизительной частотой. Считаем, что одновременно в оперативную память влезает N уникальных серчей. N
задаётся параметром. Для программы доступна безграничная файловая система.

## Start app
### build: ```make build```
### run: ```make run```

## Generate precision output for testing (with using old slow variant)
### build: ```make build```
### run: ```make run-old```

## Description
At start, I've tried to use method based on `reservoir sampling`:
- scan line by line, put values in map[query] = {sampleID, count}
- if map reached k value - delete most old query with min count 
- then compose output file: write count from queries + all discarded values

This approach works well when:
- we have little amount of queries (less than k) because we discard nothing 
- we have often occurred queries (with number << k) and other queries with 1-2 occurrences 

Problems:
- If cardinality of queries is very big we'll be unable to calculate all queries with any precision
- If some query occurred with big amount at start of file and in the end of file we will have two calculations c/2
- Big amount of duplicates

Second variant (algorithm `Count-Min sketch`):

There are many probabilistic algorithms, for example: `Count-Min Sketch`, `Count-Max Sketch`, `HyperLogLog`, `SpaceSaving` and others.
Main difference is purpose. We need to return counts not for most often queries but for every one of them, so `Count-Min sketch` was chosen.
Because it will return counts >= real values, we will not skip 1 times queries.
There still exists a problem - we need a list of queries and don't have predefined dictionary with all possible queries. If we need 
exact list of queries we can create file with name equals query in background without reading|writing to it, but it will make program slower.
For approximate list of queries `Bloom Filter` can be used (write to file only never previously encountered queries). It will cost additional memory, because additional struct will be used.
For simplicity, we can try to check if any of min_counter_table[hi] is zero (because structures are similar), but error will be more in that case (because choosing bucket algorithms are different).

This variant will store only one query per time. 
To limit memory usage we need to limit number of columns for `Count-Min Sketch` table, so k will be used as number of columns as approximate limit (where k = c * n, c = 19 * 8 / 32 = 4, len("several words query") = 19) without counting this only query.
For n = 5, k = 20 error will be = e / 20 = 0.135. If we need less error we need to take more columns.
Number of rows is fixed and equals 4 (because it will guarantee us probability = 0.96, p_error = 0.04, rows = [ln(1/p_error)]).

Possible improvements:
- Using `Bloom Filter` instead of checking  in min_counter_table to make error of false positive findings less (it will help with lost queries);
- Using `Double Hashing`, `Universal Hashing` or just faster hashes to count indexes in table (faster hashes will make algorithm faster, hash with fewer collisions will make more exact result); 
- More rows will make probability of error less;
- More columns will make more exact result.

Results:
- file with previous queries

approximate with n = 5
```
$ time ./query-counter -n 5 -ifile "./files/input.txt" -ofile "./files/output.tsv" -type approximate
2025/04/06 18:04:47 count with counter type: approximate

real    0m0.004s
user    0m0.002s
sys     0m0.002s

// output
q4	4
q1	1
q2	2
q3	5
```
approximate with n = 10 will give same results as precision
```
time ./query-counter -n 10 -ifile "./files/input.txt" -ofile "./files/output.tsv" -type approximate
2025/04/06 18:11:15 count with counter type: approximate

real    0m0.002s
user    0m0.002s
sys     0m0.001s

// output
q4	4
q1	1
q2	2
q3	3
m2	2

```

precision
```
 time ./query-counter -n 5 -ifile "./files/input.txt" -ofile "./files/exact_output.tsv" -type precision
2025/04/06 18:06:53 count with counter type: precision

real    0m0.014s
user    0m0.002s
sys     0m0.004s

// output
m2	2
q1	1
q2	2
q3	3
q4	4
```

- file with queries from task

approximate with n = 5
```
$ time ./query-counter -n 5 -ifile "./files/input_from_task.txt" -ofile "./files/output.tsv" -type approximate
2025/04/06 18:13:12 count with counter type: approximate

real    0m0.004s
user    0m0.002s
sys     0m0.003s

// output
this	2
test	2
asd	2
the	2
end	2
sad	1
is	2
only	1
```
approximate with n = 3 will give same results as precision
```
time ./query-counter -n 3 -ifile "./files/input_from_task.txt" -ofile "./files/output.tsv" -type approximate
2025/04/06 18:14:21 count with counter type: approximate

real    0m0.004s
user    0m0.000s
sys     0m0.004s


// output
this	2
test	2
asd	2
the	2
end	2
sad	1
is	1
my	1
only	1

```

precision
```
 time ./query-counter -n 5 -ifile "./files/input_from_task.txt" -ofile "./files/exact_output.tsv" -type precision
2025/04/06 18:15:55 count with counter type: precision

real    0m0.013s
user    0m0.002s
sys     0m0.004s

// output
asd	2
end	2
is	1
my	1
only	1
sad	1
test	2
the	2
this	2
```

- file with 550 rows
  approximate with n = 5
```
$ time ./query-counter -n 5 -ifile "./files/big_input.txt" -ofile "./files/output.tsv" -type approximate
2025/04/06 18:18:30 count with counter type: approximate

real    0m0.002s
user    0m0.001s
sys     0m0.002s
```


precision
```
 time ./query-counter -n 5 -ifile "./files/big_input.txt" -ofile "./files/exact_output.tsv" -type precision
2025/04/06 18:20:09 count with counter type: precision

real    0m0.029s
user    0m0.004s
sys     0m0.016s
```

With fie with 550 rows we can see real time (which includes IO operations) is much bigger for slow precision algorithm.