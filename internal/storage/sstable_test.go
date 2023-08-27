package storage

import (
	"log"
	"testing"
)

func TestPrintSegment(t *testing.T) {
	sstable := NewSSTable(&SSTableConfig{Directory: ".", FilePrefix: "A", MemtableMaxSize: 200, ValueSizeMaxBit: 8, KeySizeMaxbit: 8})
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

func TestPrintIndex(t *testing.T) {
	sstable := NewSSTable(&SSTableConfig{Directory: ".", FilePrefix: "A", MemtableMaxSize: 200, ValueSizeMaxBit: 8, KeySizeMaxbit: 8})
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

	err := sstable.IndexModel.PrintIndexFile(1)
	if err != nil {
		panic(err)
	}

}
func TestFind(t *testing.T) {
	sstable := NewSSTable(&SSTableConfig{Directory: ".", FilePrefix: "A", KeySizeMaxbit: 8, ValueSizeMaxBit: 8})
	sstable.SegmentModel.PrintSegment(1)
	sstable.IndexModel.PrintIndexFile(1)
	res, found := sstable.Find("A")
	log.Printf("find result %s", res)
	log.Println(found)
}