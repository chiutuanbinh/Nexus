package storage

import (
	"log"
	"nexus/pkg/common"
	"nexus/pkg/tree"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAVLInsert(t *testing.T) {

	sampleInput := []common.Tuple{
		{Key: "A", Value: "B"},
		{Key: "B", Value: "c"},
		{Key: "C", Value: "B"},
		{Key: "D", Value: "1"},
		{Key: "E", Value: "1"},
		{Key: "F", Value: "1"},
	}

	tree := tree.AVLTree{}
	for i := range sampleInput {
		tree.Insert(sampleInput[i].Key, sampleInput[i].Value)
		log.Println(tree.Root.ToString())
	}
	log.Println(tree.Inorder())
}

func TestAVLDelete(t *testing.T) {
	sampleInput := []common.Tuple{
		{Key: "A", Value: "B"},
		{Key: "B", Value: "c"},
		{Key: "C", Value: "B"},
		{Key: "D", Value: "1"},
		{Key: "E", Value: "1"},
		{Key: "F", Value: "1"},
	}

	tree := tree.AVLTree{}
	var err error
	for i := range sampleInput {
		err = tree.Insert(sampleInput[i].Key, sampleInput[i].Value)
		assert.Nil(t, err)
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

func TestAVLFind(t *testing.T) {

	sampleInput := []common.Tuple{
		{Key: "A", Value: "B"},
		{Key: "B", Value: "c"},
		{Key: "C", Value: "B"},
		{Key: "D", Value: "1"},
		{Key: "E", Value: "1"},
		{Key: "F", Value: "1"},
	}

	tree := tree.AVLTree{}
	for i := range sampleInput {
		tree.Insert(sampleInput[i].Key, sampleInput[i].Value)
	}
	log.Println(tree.Inorder())
	log.Println(tree.Find("A"))
	log.Println(tree.Find("F"))
	log.Println(tree.Find("E"))
}
