package common

import (
	"crypto/md5"
	"encoding/binary"

	"github.com/rs/zerolog/log"
)

type ConsistentHasher struct {
	ringSize   int
	orderedMap OrderedMap
}

// GetHashSize implements Hasher.
func (*ConsistentHasher) GetHashSize() int {
	return md5.Size
}

// Hash implements Hasher.
func (c *ConsistentHasher) FindToken(key []byte) (string, error) {
	md5Hash := md5.Sum(key)
	ringPosition := uint64(Modulo(md5Hash[:], c.ringSize))
	ringPositionArr := make([]byte, 8)
	binary.LittleEndian.PutUint64(ringPositionArr, ringPosition)
	k, nodeId, err := c.orderedMap.LowerBound(ringPositionArr)
	log.Info().Msgf("%v %vv nani", k, nodeId)
	if err != nil {
		return "", err
	}
	return string(nodeId), nil
}

func (c *ConsistentHasher) AddToken(token string) error {
	nodeIdHash := md5.Sum([]byte(token))
	ringPosition := uint64(Modulo(nodeIdHash[:], c.ringSize))
	ringPositionArr := make([]byte, 8)
	binary.LittleEndian.PutUint64(ringPositionArr, ringPosition)
	err := c.orderedMap.Put(ringPositionArr, []byte(token))
	return err
}

func (c *ConsistentHasher) DeleteToken(token string) error {
	tokenHash := md5.Sum([]byte(token))
	ringPosition := uint64(Modulo(tokenHash[:], c.ringSize))
	ringPositionArr := make([]byte, 8)
	binary.LittleEndian.PutUint64(ringPositionArr, ringPosition)
	err := c.orderedMap.Delete(ringPositionArr)
	return err
}

func CreateConsistentHasher(ringSize int) *ConsistentHasher {
	return &ConsistentHasher{
		orderedMap: CreateOrderedMap(),
		ringSize:   ringSize,
	}
}
