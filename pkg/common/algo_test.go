package common

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleArray(t *testing.T) {
	x := make([]byte, 8)
	binary.LittleEndian.PutUint64(x, 1293923981)

	y := Modulo(x, 100)
	assert.Equal(t, y, 81)

}
