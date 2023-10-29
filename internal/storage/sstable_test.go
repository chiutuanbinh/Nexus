package storage

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintSegment(t *testing.T) {
	sstable := NewSSTable(&SSTableConfig{Directory: ".", FilePrefix: "A", MemtableMaxSize: 200})
	sampleInput := []Tuple{
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
		if err != nil {
			panic(err)
		}
	}

	err := sstable.SegmentModel.PrintSegment(1)
	if err != nil {
		panic(err)
	}

}

func TestFind(t *testing.T) {
	sstable := NewSSTable(&SSTableConfig{Directory: ".", FilePrefix: "A", MemtableMaxSize: math.MaxInt32, UseHash: true})
	sampleInput := []Tuple{
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
		if err != nil {
			panic(err)
		}
	}
	sstable.Flush()
	var res string
	var found bool
	for i := range sampleInput {
		res, found = sstable.Find(sampleInput[i].Key)

		assert.True(t, found)
		assert.Equal(t, sampleInput[i].Value, res)
	}

}
