package storage

import (
	"log"
	"testing"
)

func TestFind(t *testing.T) {
	sstable := &SSTable{Directory: ".", FilePrefix: "A"}
	res, err := sstable.Find("A")
	log.Println(res)
	log.Println(err)
}
