package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistentHasher(t *testing.T) {
	c := CreateConsistentHasher(100)
	err := c.AddToken("ID")
	assert.Nil(t, err)
	nId, err := c.FindToken([]byte("G"))
	assert.Nil(t, err)
	assert.Equal(t, "ID", nId)
	err = c.AddToken("BD")
	assert.Nil(t, err)
	nId, err = c.FindToken([]byte("I"))
	assert.Nil(t, err)
	assert.Equal(t, "ID", nId)
	nId, err = c.FindToken([]byte("G"))
	assert.Nil(t, err)
	assert.Equal(t, "ID", nId)

	nId, err = c.FindToken([]byte("A"))
	assert.Nil(t, err)
	assert.Equal(t, "BD", nId)
}
