package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type SSTable struct {
	Directory        string
	LastSegmentIndex int
	FilePrefix       string
}

func (s *SSTable) Insert(tuples []Tuple) error {
	f, err := os.Create(s.getNextSegmentPath())
	if err != nil {
		return err
	}
	defer f.Close()
	for i := range tuples {
		keySize := make([]byte, 8)
		binary.LittleEndian.PutUint64(keySize, uint64(len(tuples[i].Key)))
		valueSize := make([]byte, 8)
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

func (s *SSTable) getNextSegmentPath() string {
	s.LastSegmentIndex += 1
	return s.getSegmentPath(s.LastSegmentIndex)
}

func (s *SSTable) getSegmentPath(segmentIndex int) string {
	return fmt.Sprintf("%s/%s_%v", s.Directory, s.FilePrefix, segmentIndex)
}

func (s *SSTable) Find(key string) (string, error) {

	for i := s.LastSegmentIndex; i > 0; i-- {
		f, err := os.Open(s.getSegmentPath(i))
		if err != nil {
			return "", err
		}
		defer f.Close()
		sizeBuf := make([]byte, 8)
		dataBuf := make([]byte, 0)
		for {
			_, err := f.Read(sizeBuf)
			if err == io.EOF {
				return "", nil
			}
			size := binary.LittleEndian.Uint64(sizeBuf)
			dataBuf = dataBuf[:size]
			_, err = f.Read(dataBuf)
			data := string(dataBuf[:])
			fmt.Println(data)
		}
	}
	return "", nil
}

func (s *SSTable) PrintTheFirstSegment() error {
	f, err := os.Open(s.getSegmentPath(0))
	if err != nil {
		return err
	}
	defer f.Close()
	sizeBuf := make([]byte, 8)

	for {
		_, err := f.Read(sizeBuf)
		if err == io.EOF {
			return nil
		}
		size := binary.LittleEndian.Uint64(sizeBuf)
		dataBuf := make([]byte, size)
		_, err = f.Read(dataBuf)
		data := string(dataBuf[:])
		fmt.Println(data)
	}
}
