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
		key, err := readNextValue(f, s.Config.KeySizeMaxbit)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		value, err := readNextValue(f, s.Config.ValueSizeMaxBit)
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
	log.Printf("SIZE %v", size)
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
		keySize := make([]byte, s.Config.KeySizeMaxbit)
		binary.LittleEndian.PutUint64(keySize, uint64(16))
		valueSize := make([]byte, s.Config.ValueSizeMaxBit)
		binary.LittleEndian.PutUint64(valueSize, uint64(len(tuples[i].Value)))
		_, err := f.Write(keySize)
		if err != nil {
			return err
		}
		keyByteArr, err := hex.DecodeString(tuples[i].Key)
		if err != nil {
			return err
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
	keySizeArr := make([]byte, s.Config.KeySizeMaxbit)
	_, err = f.Read(keySizeArr)
	if err != nil {
		panic(err)
	}

	keySize := binary.LittleEndian.Uint64(keySizeArr)
	keyByteArr := make([]byte, keySize)
	log.Printf("KEY SIZE %v %v", keySize, s.Config.KeySizeMaxbit)
	_, err = f.Read(keyByteArr)
	if err != nil {
		panic(err)
	}

	key := hex.EncodeToString(keyByteArr)
	log.Printf("KEY %v", key)
	valueSizeArr := make([]byte, s.Config.ValueSizeMaxBit)

	_, err = f.Read(valueSizeArr)
	if err != nil {
		panic(err)
	}
	valueSize := binary.LittleEndian.Uint64(valueSizeArr)
	log.Printf("Value size %v", valueSize)
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
		keyByteArr, err := hex.DecodeString(tuples[index].Key)
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
		nextPos += i.Config.KeySizeMaxbit + i.Config.KeySize + +i.Config.ValueSizeMaxBit + valueSize

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

func (i *IndexFileModel) Find(key string) (int, int) {
	keyHex, err := hex.DecodeString(key)
	if err != nil {
		panic(fmt.Sprintf("Key is not hashed %v", key))
	}
	// step := i.Config.KeySize + 8
	for index := i.LastIndex; index > 0; index-- {

		f, err := os.Open(i.getIndexPath(index))
		if err != nil {
			panic(fmt.Sprintf("Index file %v not found", i.getIndexPath(index)))
		}
		defer f.Close()

		for {
			keyBuf := make([]byte, i.Config.KeySize)
			_, err = f.Read(keyBuf)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			posBuf := make([]byte, 8)
			_, err = f.Read(posBuf)
			if bytes.Equal(keyHex, keyBuf) {
				return index, int(binary.LittleEndian.Uint64(posBuf))
			}

		}
	}
	return -1, -1
}

func (i *IndexFileModel) keyCount(fileSize int) int {
	panic("IndexFileModel keyCount not implemented")
}

func (i *IndexFileModel) PrintIndexFile(index int) error {
	f, err := os.Open(i.getIndexPath(index))
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		keyBuf := make([]byte, i.Config.KeySize)
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
	KeySize          int
	KeySizeMaxbit    int
	ValueSizeMaxBit  int
	MemtableMaxSize  int
}

type SSTable struct {
	Config       *SSTableConfig
	SegmentModel *SegmentFileModel
	IndexModel   *IndexFileModel
	Memtable     Memtable
	Hasher       hashing.Hasher
}

func NewSSTable(config *SSTableConfig) *SSTable {
	var hasher hashing.Hasher = &hashing.MD5Hasher{}
	config.KeySize = hasher.GetHashSize()
	return &SSTable{
		Config: config,
		SegmentModel: &SegmentFileModel{
			Config: config,
		},
		IndexModel: &IndexFileModel{
			Config:    config,
			LastIndex: 1,
		},
		Hasher:   hasher,
		Memtable: &AVLTree{},
	}
}

func (s *SSTable) Insert(key string, value string) error {
	hashedKey := s.Hasher.Hash(key)
	log.Println(hashedKey)
	err := s.Memtable.Insert(string(hashedKey), value)
	if err != nil {
		return err
	}
	if s.Memtable.Size() >= s.Config.MemtableMaxSize {
		//TODO: Make this run on separated routine
		return s.Flush()
	}
	return nil
}

func (s *SSTable) Find(key string) (string, bool) {
	hashedKey := s.Hasher.Hash(key)
	v, ok := s.Memtable.Find(hashedKey)
	if ok {
		return v, true
	}
	fileIndex, pos := s.IndexModel.Find(hashedKey)
	log.Printf("%v %v %X\n", hashedKey, fileIndex, pos)
	if fileIndex == -1 && pos == -1 {
		return "", false
	}
	t, err := s.SegmentModel.Get(fileIndex, pos)
	if err != nil {
		panic(err)
	}
	return t.Value, true
}

func (s *SSTable) Flush() error {
	tuples := s.Memtable.List()
	err := s.IndexModel.Flush(tuples)
	if err != nil {
		return err
	}
	err = s.SegmentModel.Flush(tuples)
	if err != nil {
		return err
	}
	err = s.Memtable.Clear()
	return err
}

func (s *SSTable) Cleanup() {
	// TODO: handle clean up on sigint or sigterm
}
