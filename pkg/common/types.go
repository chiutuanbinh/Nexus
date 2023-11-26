package common

import "fmt"

type Tuple struct {
	Key   []byte
	Value []byte
}

type Hasher interface {
	Hash(in []byte) []byte
	GetHashSize() int
}

type NotFoundError struct {
	Value string
}

// Error implements error.
func (n NotFoundError) Error() string {
	return fmt.Sprintf("Value %v cannot be found", n.Value)
}

var _ error = NotFoundError{}
