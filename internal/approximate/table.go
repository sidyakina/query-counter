package approximate

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
)

type Table struct {
	columns     int
	hashToCount [][]int32
}

const (
	rows = 4
	c    = 4
)

func NewTable(n int) *Table {
	columns := c * n // additional info in README.md

	hashToCount := make([][]int32, rows)
	for i := 0; i < rows; i++ {
		hashToCount[i] = make([]int32, columns)
	}

	return &Table{columns: columns, hashToCount: hashToCount}
}

func (t *Table) getIndexes(query string) [4]int {
	indexes := [4]int{}
	indexes[0] = int(binary.LittleEndian.Uint64(md5.New().Sum([]byte(query))) % uint64(t.columns))
	indexes[1] = int(binary.LittleEndian.Uint64(sha1.New().Sum([]byte(query))) % uint64(t.columns))
	indexes[2] = int(binary.LittleEndian.Uint64(sha256.New().Sum([]byte(query))) % uint64(t.columns))
	indexes[3] = int(binary.LittleEndian.Uint64(sha512.New().Sum([]byte(query))) % uint64(t.columns))

	return indexes
}

func (t *Table) Add(query string) (exists bool) {
	exists = true

	indexes := t.getIndexes(query)
	for i := 0; i < rows; i++ {
		idx := indexes[i]

		v := t.hashToCount[i][idx]
		if v == 0 {
			exists = false
		}

		t.hashToCount[i][idx] = v + 1
	}

	return exists
}

func (t *Table) Count(query string) int32 {
	count := int32(-1)

	indexes := t.getIndexes(query)
	for i := 0; i < rows; i++ {
		idx := indexes[i]

		v := t.hashToCount[i][idx]
		if count == -1 || count > v {
			count = v
		}
	}

	return count

}
