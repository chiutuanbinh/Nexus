package storage

import "nexus/pkg/common"

type Memtable interface {
	Insert(key string, value string) error
	Delete(key string) (bool, error)
	Find(key string) (string, bool)
	Clear() error
	List() []common.Tuple
	Size() int
}

// If there are 2 equal key, replace the value, we do not allow duplicate key
