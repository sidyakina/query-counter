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
- after read each line compose output file: write count from queries + all discarded values
This approach works well when:
- we have little amount of queries (less than k) because we discard nothing 
- we have often occurred queries (with number << k) and other queries with 1-2 occurrences 
Problems:
- If cardinality of queries very big we'll be unable to calculate all queries with any precision
- If some query occurred with big amount at start of file and in the end of file we will have two calculation c/2
- Big amount of duplicates

Second variant (algorithm `Count-Min sketch`):
There are many probabilistic algorithms, for example: `Count-Min Sketch`, `Count-Max Sketch`, `HyperLogLog`, `SpaceSaving` and others.
Main difference is purpose. We need to return counts not for most often queries but for every one of them, so `Count-Min sketch` was chosen.
Because it will return counts >= real values, so we will not skip 1 times queries.
There still exists a problem - we need a list and don't have predefined dictionary with all possible queries. If we need 
exact list of queries we can create file with name in background without reading|writing to it, but it will make program more slow.
For approximate list of queries `Bloom Filter` can be used (write to file only never encountered queries). It will cost additional memory, because additional struct will be used.
For simplicity, we can try to check if any of min_counter_table[hi] is zero (because structures are similar), but error will be more in that case (because choosing bucket algorithms are different).

This variant will be store only one query per time. 
To limit memory usage we need to limit number of columns for `Count-Min Sketch` table, so k will be used as number of columns as approximate limit (where k = c * n, c = 19 * 8 / 32 = 4, len("several words query") = 19) without counting this only query.
For n = 5, k = 20 error will be = e / 20 = 0.135. If we need less error we need to take more columns.
Number of rows is fixed and equals 4 (because it will guarantee us probability = 0.96, p_error = 0.04, rows = [ln(1/p_error)]).

Possible improvements:
- Using `Bloom Filter` instead of checking  in min_counter_table to make error of false positive findings less (it will help with lost queries);
- Using `Double Hashing`, `Universal Hashing` or just more fast hashes to count indexes in table (more fast will make algorithm faster, hash with fewer collisions will make more exact result); 
- More rows will make probability of error less;
- More columns (with more memory) will make more exact result.