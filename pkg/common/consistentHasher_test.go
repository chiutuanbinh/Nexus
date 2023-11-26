package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistentHasher(t *testing.T) {
	c := CreateConsistentHasher(100)
	err := c.AddNode("ID")
	assert.Nil(t, err)
	nId, err := c.FindNode([]byte("G"))
	assert.Nil(t, err)
	assert.Equal(t, "ID", nId)
	err = c.AddNode("BD")
	assert.Nil(t, err)
	nId, err = c.FindNode([]byte("I"))
	assert.Nil(t, err)
	assert.Equal(t, "ID", nId)
	nId, err = c.FindNode([]byte("G"))
	assert.Nil(t, err)
	assert.Equal(t, "ID", nId)

	nId, err = c.FindNode([]byte("A"))
	assert.Nil(t, err)
	assert.Equal(t, "BD", nId)
}
