package hash

import (
	"crypto/md5"
	"encoding/hex"
)

type MD5Hasher struct {
}

func (md5Hasher MD5Hasher) Hash(in string) string {
	res := md5.Sum([]byte(in))
	return hex.EncodeToString(res[:])
}

func (md5Hasher MD5Hasher) GetHashSize() int {
	return 16
}
