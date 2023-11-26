package common

import (
	"crypto/md5"
)

type MD5Hasher struct {
}

var _ Hasher = &MD5Hasher{}

func (md5Hasher MD5Hasher) Hash(in []byte) []byte {
	res := md5.Sum(in)
	return res[:]
}

func (md5Hasher MD5Hasher) GetHashSize() int {
	return md5.Size
}
