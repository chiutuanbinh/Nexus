package storage

import (
	"bytes"
	"fmt"
	"nexus/pkg/common"
)

type SSTableConfig struct {
	Directory        string
	FilePrefix       string
	SegmentThreshold uint32
	MemtableMaxSize  int
	UseHash          bool
}

type SSTable struct {
	Config        *SSTableConfig
	segmentModel  *SegmentFileModel
	indexModel    *IndexFileModel
	memtable      Memtable
	hasher        common.Hasher
	flushCallBack func() error
}

func NewSSTable(config *SSTableConfig, flushCallBack func() error) *SSTable {
	hasher := &common.MD5Hasher{}

	sstable := &SSTable{
		Config: config,
		segmentModel: CreateSegmentFileModel(
			&SegmentFileModelConfig{
				Directory:  config.Directory,
				FilePrefix: config.FilePrefix,
			},
		),
		indexModel: CreateIndexFileModel(
			&IndexFileModelConfig{
				Directory:  config.Directory,
				FilePrefix: config.FilePrefix,
			},
		),
		hasher:        hasher,
		memtable:      common.CreateAVLTree(),
		flushCallBack: flushCallBack,
	}

	return sstable
}

func (s *SSTable) preprocess(v string) ([]byte, error) {
	if s.Config.UseHash {
		keyBytes := s.hasher.Hash([]byte(v))
		return keyBytes, nil
	} else {
		if len(v) > KEY_SIZE_IN_BYTE {
			return nil, fmt.Errorf("Key %s larger than max size %v", v, KEY_SIZE_IN_BYTE)
		}
		vArr := []byte(v)
		padding := make([]byte, KEY_SIZE_IN_BYTE-len(vArr))
		vArr = append(vArr, padding...)
		return vArr, nil
	}
}

func (s *SSTable) Insert(key string, value string) error {
	preprocessedKey, err := s.preprocess(key)
	if err != nil {
		return err
	}
	err = s.memtable.Insert(preprocessedKey, []byte(value))
	if err != nil {
		return err
	}
	if s.memtable.Size() >= s.Config.MemtableMaxSize {
		//TODO: Make this run on separated routine
		return s.Flush()
	}
	return nil
}

func (s *SSTable) Delete(key string) error {
	return s.Insert(key, string(getTombStone()))
}

func (s *SSTable) Find(key string) (string, bool) {
	// var cost = time.Now().UnixMilli()
	preprocessedKey, err := s.preprocess(key)
	if err != nil {
		return "", false
	}

	// s1 := time.Now().UnixMilli() - cost
	v, err := s.memtable.Find(preprocessedKey)
	// log.Debug().Msgf("Find key %v processed %v value %v", key, string(preprocessedKey), v)
	if err == nil {
		if bytes.Equal(v, getTombStone()) {
			return "", false
		}
		return string(v), true
	}
	// s2 := time.Now().UnixMilli() - cost
	fileIndex, pos := s.indexModel.Find(preprocessedKey)
	if fileIndex == -1 && pos == -1 {
		return "", false
	}
	// s3 := time.Now().UnixMilli() - cost
	t, err := s.segmentModel.Get(fileIndex, pos)
	if err != nil {
		panic(err)
	}
	// s4 := time.Now().UnixMilli() - cost
	if bytes.Equal(t.Value, getTombStone()) {
		return "", false
	}
	// log.Info().Msgf("Cost of finding a value: preprocess %v mem %v index %v seg %v", s1, s2, s3, s4)
	return string(t.Value), true
}

func (s *SSTable) Flush() error {
	tuples := s.memtable.List()
	err := s.indexModel.Flush(tuples)
	if err != nil {
		return fmt.Errorf("Index file flush failed %w", err)
	}
	err = s.segmentModel.Flush(tuples)
	if err != nil {
		return fmt.Errorf("Segment file flush failed %w", err)
	}

	err = s.memtable.Clear()
	if err != nil {
		return fmt.Errorf("Memtable clear failed %w", err)
	}
	err = s.flushCallBack()
	return err
}
