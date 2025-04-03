package approximate

import "fmt"

type QueryInfo struct {
	count    int
	lastSeen int
}

type Queries struct {
	queries  map[string]*QueryInfo
	capacity int
}

func NewQueries(capacity int) *Queries {
	return &Queries{capacity: capacity, queries: make(map[string]*QueryInfo)}
}

func (q *Queries) Add(query string, sampleNumber int) (discardedQuery string, count int) {
	info, ok := q.queries[query]
	if ok {
		info.count = info.count + 1
		info.lastSeen = sampleNumber
	} else {
		q.queries[query] = &QueryInfo{count: 1, lastSeen: sampleNumber}
	}

	if len(q.queries) > q.capacity {
		fmt.Printf("add query: %v with %v\n", query, sampleNumber)
		return q.Discard()
	}

	return "", 0
}

func (q *Queries) Discard() (discardedQuery string, count int) {
	lastSeen := -1

	for k, v := range q.queries {
		fmt.Printf("checking: %v - %v, %v; ", k, v.lastSeen, v.count)
		// remove most old query with min count
		candidateToRemove := lastSeen == -1 || (lastSeen == v.lastSeen && count > v.count) || lastSeen > v.lastSeen
		if candidateToRemove {
			fmt.Print("it's candidate")
			lastSeen = v.lastSeen
			discardedQuery = k
			count = v.count
		}
		fmt.Println()
	}

	delete(q.queries, discardedQuery)

	return discardedQuery, count
}
