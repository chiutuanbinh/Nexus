package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"nexus/internal/hashing"
	"os"
)

type SegmentFileModel struct {
	Config           *SSTableConfig
	LastSegmentIndex int
}

func NewSegmentFileModel() *SegmentFileModel {
	return nil
}

func (s *SegmentFileModel) getNextSegmentPath() string {
	s.LastSegmentIndex += 1
	return s.getSegmentPath(s.LastSegmentIndex)
}

func (s *SegmentFileModel) getSegmentPath(segmentIndex int) string {
	return fmt.Sprintf("%s/%s_%v.seg", s.Config.Directory, s.Config.FilePrefix, segmentIndex)
}

func (s *SegmentFileModel) PrintSegment(segmentIndex int) error {
	f, err := os.Open(s.getSegmentPath(segmentIndex))
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		key, err := readNextValue(f, KEY_SIZE_MAX_BYTE_LENGTH)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		value, err := readNextValue(f, VALUE_SIZE_MAX_BYTE_LENGTH)
		if err != nil {
			return err
		}
		log.Printf("%v:%v", key, value)
	}
}

func readNextValue(f *os.File, sizeBufMaxBit int) (string, error) {
	sizeBuf := make([]byte, sizeBufMaxBit)
	_, err := f.Read(sizeBuf)
	if err == io.EOF {
		return "", err
	}
	size := binary.LittleEndian.Uint64(sizeBuf)
	dataBuf := make([]byte, size)
	_, err = f.Read(dataBuf)
	if err == io.EOF {
		return "", err
	}
	return string(dataBuf[:]), nil
}

func (s *SegmentFileModel) Flush(tuples []Tuple) error {
	f, err := os.Create(s.getNextSegmentPath())
	if err != nil {
		return err
	}
	defer f.Close()
	for i := range tuples {
		keySize := make([]byte, KEY_SIZE_MAX_BYTE_LENGTH)
		//TODO: replace the size with configurable value
		binary.LittleEndian.PutUint64(keySize, uint64(16))
		valueSize := make([]byte, VALUE_SIZE_MAX_BYTE_LENGTH)
		binary.LittleEndian.PutUint64(valueSize, uint64(len(tuples[i].Value)))
		_, err := f.Write(keySize)
		if err != nil {
			return err
		}
		var keyByteArr []byte
		if s.Config.UseHash {
			keyByteArr, err = hex.DecodeString(tuples[i].Key)
			if err != nil {
				return err
			}
		} else {
			keyByteArr = []byte(tuples[i].Key)
		}

		_, err = f.Write(keyByteArr)
		if err != nil {
			return err
		}
		_, err = f.Write(valueSize)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(tuples[i].Value))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SegmentFileModel) Get(fileIndex int, pos int) (Tuple, error) {
	f, err := os.Open(s.getSegmentPath(fileIndex))
	if err != nil {
		panic(fmt.Sprintf("Segment file %v not found", s.getSegmentPath(fileIndex)))
	}
	defer f.Close()

	_, err = f.Seek(int64(pos), 0)
	if err != nil {
		panic(err)
	}
	keySizeArr := make([]byte, KEY_SIZE_MAX_BYTE_LENGTH)
	_, err = f.Read(keySizeArr)
	if err != nil {
		panic(err)
	}

	keySize := binary.LittleEndian.Uint64(keySizeArr)
	keyByteArr := make([]byte, keySize)
	_, err = f.Read(keyByteArr)
	if err != nil {
		panic(err)
	}

	key := hex.EncodeToString(keyByteArr)
	valueSizeArr := make([]byte, VALUE_SIZE_MAX_BYTE_LENGTH)

	_, err = f.Read(valueSizeArr)
	if err != nil {
		panic(err)
	}
	valueSize := binary.LittleEndian.Uint64(valueSizeArr)
	valueArr := make([]byte, valueSize)
	_, err = f.Read(valueArr)
	if err != nil {
		panic(err)
	}
	value := string(valueArr)
	return Tuple{Key: key, Value: value}, nil

}

// The index file contain fixed size entry, support fast binary search for the correct key position
type IndexFileModel struct {
	Config    *SSTableConfig
	LastIndex int
}

func (i *IndexFileModel) Flush(tuples []Tuple) error {
	f, err := os.Create(i.getNextIndexPath())
	if err != nil {
		return err
	}
	defer f.Close()
	nextPos := 0
	for index := range tuples {
		var keyByteArr []byte
		if i.Config.UseHash {
			keyByteArr, err = hex.DecodeString(tuples[index].Key)
		} else {
			keyByteArr = []byte(tuples[index].Key)[:KEY_SIZE_IN_BYTE]
		}

		if err != nil {
			return err
		}
		_, err = f.Write(keyByteArr)
		if err != nil {
			return err
		}
		nextPosByteArr := make([]byte, 8)
		binary.LittleEndian.PutUint64(nextPosByteArr, uint64(nextPos))
		_, err = f.Write(nextPosByteArr)
		if err != nil {
			return err
		}
		valueSize := len(tuples[index].Value)
		nextPos += KEY_SIZE_MAX_BYTE_LENGTH + KEY_SIZE_IN_BYTE + VALUE_SIZE_MAX_BYTE_LENGTH + valueSize

	}
	return nil
}

func (i *IndexFileModel) getNextIndexPath() string {
	i.LastIndex += 1
	return i.getIndexPath(i.LastIndex)
}

func (i *IndexFileModel) getIndexPath(nextIndex int) string {
	return fmt.Sprintf("%s/%s_%v.index", i.Config.Directory, i.Config.FilePrefix, nextIndex)
}

func (i *IndexFileModel) Find(key []byte) (int, int) {
	for index := i.LastIndex; index > 0; index-- {
		f, err := os.Open(i.getIndexPath(index))
		if err != nil {
			panic(fmt.Sprintf("Index file %v not found", i.getIndexPath(index)))
		}
		defer f.Close()

		for {
			keyBuf := make([]byte, KEY_SIZE_IN_BYTE)
			_, err = f.Read(keyBuf)
			if err == io.EOF {
				break
			}
			// log.Printf("%v %v", key, keyBuf)
			if err != nil {
				panic(err)
			}
			posBuf := make([]byte, POSITION_OFFSET_SIZE_IN_BYTE)
			_, err = f.Read(posBuf)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			if bytes.Equal(key, keyBuf) {
				return index, int(binary.LittleEndian.Uint64(posBuf))
			}

		}
	}
	return -1, -1
}

func (i *IndexFileModel) KeyCount(fileSize int) int {
	panic("IndexFileModel keyCount not implemented")
}

func (i *IndexFileModel) PrintIndexFile(index int) error {
	f, err := os.Open(i.getIndexPath(index))
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		keyBuf := make([]byte, KEY_SIZE_IN_BYTE)
		_, err := f.Read(keyBuf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		posBuf := make([]byte, 8)
		_, err = f.Read(posBuf)
		if err != nil {
			return err
		}
		pos := binary.LittleEndian.Uint64(posBuf)
		log.Printf("%v pos %X\n", hex.EncodeToString(keyBuf), pos)

	}
}

// TODO: make these config available when reading any segment file
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
	hasher        hashing.Hasher
	flushCallBack func() error
}

func NewSSTable(config *SSTableConfig, flushCallBack func() error) *SSTable {
	var hasher hashing.Hasher = &hashing.MD5Hasher{}

	sstable := &SSTable{
		Config: config,
		segmentModel: &SegmentFileModel{
			Config: config,
		},
		indexModel: &IndexFileModel{
			Config:    config,
			LastIndex: 0,
		},
		hasher:        hasher,
		memtable:      &AVLTree{},
		flushCallBack: flushCallBack,
	}

	return sstable
}

func (s *SSTable) preprocess(v string) ([]byte, error) {
	if s.Config.UseHash {
		hashedKey := s.hasher.Hash(v)
		keyBytes, err := hex.DecodeString(hashedKey)
		if err != nil {
			return nil, err
		}
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
	err = s.memtable.Insert(string(preprocessedKey), value)
	if err != nil {
		return err
	}
	if s.memtable.Size() >= s.Config.MemtableMaxSize {
		//TODO: Make this run on separated routine
		return s.Flush()
	}
	return nil
}

func (s *SSTable) Find(key string) (string, bool) {
	preprocessedKey, err := s.preprocess(key)
	if err != nil {
		return "", false
	}
	v, ok := s.memtable.Find(string(preprocessedKey))
	if ok {
		return v, true
	}
	fileIndex, pos := s.indexModel.Find(preprocessedKey)
	if fileIndex == -1 && pos == -1 {
		return "", false
	}
	t, err := s.segmentModel.Get(fileIndex, pos)
	if err != nil {
		panic(err)
	}
	return t.Value, true
}

func (s *SSTable) Flush() error {
	tuples := s.memtable.List()
	err := s.indexModel.Flush(tuples)
	if err != nil {
		return err
	}
	err = s.segmentModel.Flush(tuples)
	if err != nil {
		return err
	}

	err = s.memtable.Clear()
	if err != nil {
		return err
	}
	err = s.flushCallBack()
	return err
}
