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
		log.Println(tree.Root.ToString())
	}

}

func TestAVLDelete(t *testing.T) {
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
	}

	log.Println(tree.Root.ToString())
	log.Println("DELETE A")
	tree.Delete("A")
	log.Println(tree.Root.ToString())
	log.Println("DELETE B")
	tree.Delete("B")
	log.Println(tree.Root.ToString())
	log.Println("DELETE D")
	tree.Delete("D")
	log.Println(tree.Root.ToString())
	log.Println("DELETE D")
	tree.Delete("D")
	log.Println(tree.Root.ToString())
	log.Println("DELETE E")
	tree.Delete("E")
	log.Println(tree.Root.ToString())
}
