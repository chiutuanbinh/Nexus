package storage

import (
	"log"
	"testing"
)

func TestAVLInsertt(t *testing.T) {

	sampleInput := []struct {
		key   string
		value string
	}{
		{key: "A", value: "B"},
		{key: "B", value: "c"},
		{key: "C", value: "B"},
		{key: "D", value: "1"},
		{key: "E", value: "1"},
		{key: "F", value: "1"},
	}

	tree := AVLTree{}
	for i := range sampleInput {
		tree.Insert(sampleInput[i].key, sampleInput[i].value)
		log.Println(i)
	}
	log.Println(tree.Root.ToString())
}
