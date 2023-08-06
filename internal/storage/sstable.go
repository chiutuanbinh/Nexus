package storage

import (
	"encoding/binary"
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
	return fmt.Sprintf("%s/%s_seg_%v", s.Config.Directory, s.Config.FilePrefix, segmentIndex)
}

func (s *SegmentFileModel) PrintSegment(segmentIndex int) error {
	f, err := os.Open(s.getSegmentPath(segmentIndex))
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		key, err := s.readNextValue(f)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		value, err := s.readNextValue(f)
		if err != nil {
			return err
		}
		log.Printf("%v:%v", key, value)
	}
}

func (s *SegmentFileModel) readNextValue(f *os.File) (string, error) {
	sizeBuf := make([]byte, 8)
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
		keySize := make([]byte, s.Config.KeySizeMaxbit)
		binary.LittleEndian.PutUint64(keySize, uint64(len(tuples[i].Key)))
		valueSize := make([]byte, s.Config.ValueSizeMaxBit)
		binary.LittleEndian.PutUint64(valueSize, uint64(len(tuples[i].Value)))

		_, err := f.Write(keySize)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(tuples[i].Key))
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

// The index file contain fixed size entry, support fast binary search for the correct key
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
		_, err := f.Write([]byte(tuples[index].Key))
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
		nextPos += i.Config.KeySizeMaxbit + valueSize

	}
	return nil
}

func (i *IndexFileModel) getNextIndexPath() string {
	i.LastIndex += 1
	return i.getIndexPath(i.LastIndex)
}

func (i *IndexFileModel) getIndexPath(nextIndex int) string {
	return fmt.Sprintf("%s/%s_index_%v", i.Config.Directory, i.Config.FilePrefix, nextIndex)
}

func (i *IndexFileModel) Find(key string) Tuple {
	return Tuple{}
}

// TODO: make these config available when reading any segment file
type SSTableConfig struct {
	Directory        string
	FilePrefix       string
	SegmentThreshold uint32
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
	return &SSTable{
		Config: config,
		SegmentModel: &SegmentFileModel{
			Config: config,
		},
		IndexModel: &IndexFileModel{
			Config: config,
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
		tuples := s.Memtable.List()
		err = s.IndexModel.Flush(tuples)
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
	return nil
}

func (s *SSTable) Find(key string) (string, error) {
	return "", nil
}
