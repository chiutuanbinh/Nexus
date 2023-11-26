package storage

import "nexus/pkg/common"

type Memtable interface {
	common.AVLTree
}

// If there are 2 equal key, replace the value, we do not allow duplicate key
