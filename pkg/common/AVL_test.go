package common

import (
	"log"
	"testing"
)

func TestAVLInsert(t *testing.T) {

	sampleInput := []struct {
		Key   string
		Value string
	}{
		{Key: "A", Value: "B"},
		{Key: "B", Value: "c"},
		{Key: "C", Value: "B"},
		{Key: "D", Value: "1"},
		{Key: "E", Value: "1"},
		{Key: "F", Value: "1"},
	}

	tree := CreateAVLTree()
	for i := range sampleInput {
		tree.Insert([]byte(sampleInput[i].Key), []byte(sampleInput[i].Value))
	}
	log.Println(tree.List())
}

func TestAVLFind(t *testing.T) {

	sampleInput := []struct {
		Key   string
		Value string
	}{
		{Key: "A", Value: "B"},
		{Key: "B", Value: "c"},
		{Key: "C", Value: "B"},
		{Key: "D", Value: "1"},
		{Key: "E", Value: "1"},
		{Key: "F", Value: "1"},
	}

	tree := CreateAVLTree()
	for i := range sampleInput {
		tree.Insert([]byte(sampleInput[i].Key), []byte(sampleInput[i].Value))
	}
	log.Println(tree.List())
	log.Println(tree.Find([]byte("A")))
	log.Println(tree.Find([]byte("F")))
	log.Println(tree.Find([]byte("E")))
}
