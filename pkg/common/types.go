package common

import "fmt"

type BstTree interface {
	Insert(key []byte, value []byte) error
	Delete(key []byte) error
	Find(key []byte) ([]byte, error)
	LowerBound(key []byte) ([]byte, []byte, error)
	UpperBound(key []byte) ([]byte, []byte, error)
	Clear() error
	Size() int
	List() []Tuple
	NodeCount() int
}
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
