package hashing

import (
	"crypto/md5"
	"encoding/hex"
)

type Hasher interface {
	Hash(in string) string
	GetHashSize() int
}

type MD5Hasher struct {
}

func (md5Hasher MD5Hasher) Hash(in string) string {
	res := md5.Sum([]byte(in))
	return hex.EncodeToString(res[:])
}

func (md5Hasher MD5Hasher) GetHashSize() int {
	return 16
}
