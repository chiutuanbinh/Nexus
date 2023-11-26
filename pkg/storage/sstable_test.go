package storage

import (
	"fmt"
	"math"

	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type T struct {
	Key   string
	Value string
}

func TestFind(t *testing.T) {
	sstable := NewSSTable(&SSTableConfig{Directory: ".", FilePrefix: "A", MemtableMaxSize: math.MaxInt32, UseHash: true}, func() error { return nil })
	sampleInput := []T{
		{Key: "A", Value: "B"},
		{Key: "B", Value: "c"},
		{Key: "C", Value: "B"},
		{Key: "D", Value: "1"},
		{Key: "E", Value: "1"},
		{Key: "F", Value: "1"},
		{Key: "G", Value: "B"},
		{Key: "H", Value: "1"},
		{Key: "I", Value: "1"},
		{Key: "J", Value: "1"},
	}
	for i := range sampleInput {
		err := sstable.Insert(sampleInput[i].Key, sampleInput[i].Value)
		assert.Nil(t, err)
	}
	err := sstable.Flush()
	assert.Nil(t, err)
	var res string
	var found bool
	for i := range sampleInput {
		res, found = sstable.Find(sampleInput[i].Key)

		assert.True(t, found)
		assert.Equal(t, sampleInput[i].Value, res)
	}

	err = sstable.Delete("A")
	assert.Nil(t, err)
	log.Debug().Msg(fmt.Sprintf("%v", sstable.memtable.List()))
	err = sstable.Flush()
	assert.Nil(t, err)
	res, found = sstable.Find("A")

	assert.False(t, found)
	assert.Empty(t, res)
}
