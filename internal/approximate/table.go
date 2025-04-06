package approximate

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"hash"
)

const (
	rows = 4
	c    = 4
)

type Table struct {
	columns     int
	hashToCount [rows][]int32
}

func NewTable(n int) *Table {
	columns := c * n // additional info in README.md

	hashToCount := [rows][]int32{}
	for i := 0; i < rows; i++ {
		hashToCount[i] = make([]int32, columns)
	}

	return &Table{columns: columns, hashToCount: hashToCount}
}

func (t *Table) getIndexes(query string) [rows]int {
	indexes := [rows]int{}

	for i, hashFunc := range [rows]hash.Hash{md5.New(), sha1.New(), sha256.New(), sha512.New()} {
		sum := hashFunc.Sum([]byte(query))
		indexes[i] = int(binary.LittleEndian.Uint64(sum) % uint64(t.columns))
	}

	return indexes
}

func (t *Table) Add(query string) (exists bool) {
	exists = true

	indexes := t.getIndexes(query)
	for i := 0; i < rows; i++ {
		index := indexes[i]

		count := t.hashToCount[i][index]
		// count is zero - we never encountered query with such hash index
		// but if all count is not zero it can also mean we have encountered queries with collision indexes
		// query excluding from dictionary will be false positive
		if count == 0 {
			exists = false
		}

		t.hashToCount[i][index] = count + 1
	}

	return exists
}

func (t *Table) Count(query string) int32 {
	result := int32(-1)

	indexes := t.getIndexes(query)
	for i := 0; i < rows; i++ {
		index := indexes[i]

		count := t.hashToCount[i][index]
		if result == -1 || result > count {
			result = count
		}
	}

	return result

}
