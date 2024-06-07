package storage

import "nexus/pkg/common"

type Memtable interface {
	common.BstTree
}

// If there are 2 equal key, replace the value, we do not allow duplicate key
