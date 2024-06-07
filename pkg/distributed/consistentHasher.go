package distributed

import (
	"crypto/md5"
	"encoding/binary"
	"nexus/pkg/common"
)

type ConsistentHasher struct {
	ringSize      int
	tokenPosition common.OrderedMap
}

func (*ConsistentHasher) GetHashSize() int {
	return md5.Size
}

// Return the token this key is mapped to
func (c *ConsistentHasher) FindToken(key []byte) (string, error) {
	md5Hash := md5.Sum(key)
	ringPosition := uint64(common.Modulo(md5Hash[:], c.ringSize))
	ringPositionArr := make([]byte, 8)
	binary.LittleEndian.PutUint64(ringPositionArr, ringPosition)
	_, nodeId, err := c.tokenPosition.LowerBound(ringPositionArr)
	if err != nil {
		return "", err
	}
	return string(nodeId), nil
}

func (c *ConsistentHasher) AddToken(token string) error {
	nodeIdHash := md5.Sum([]byte(token))
	ringPosition := uint64(common.Modulo(nodeIdHash[:], c.ringSize))
	ringPositionArr := make([]byte, 8)
	binary.LittleEndian.PutUint64(ringPositionArr, ringPosition)
	err := c.tokenPosition.Put(ringPositionArr, []byte(token))
	return err
}

func (c *ConsistentHasher) DeleteToken(token string) error {
	tokenHash := md5.Sum([]byte(token))
	ringPosition := uint64(common.Modulo(tokenHash[:], c.ringSize))
	ringPositionArr := make([]byte, 8)
	binary.LittleEndian.PutUint64(ringPositionArr, ringPosition)
	err := c.tokenPosition.Delete(ringPositionArr)
	return err
}

func CreateConsistentHasher(ringSize int) *ConsistentHasher {
	return &ConsistentHasher{
		tokenPosition: common.CreateOrderedMap(),
		ringSize:      ringSize,
	}
}
