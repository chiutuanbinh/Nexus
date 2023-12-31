package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"nexus/pkg/common"
	"os"
	"regexp"
	"strconv"

	"github.com/rs/zerolog/log"
)

type SegmentFileModelConfig struct {
	Directory  string
	FilePrefix string
}

type SegmentFileModel struct {
	Config           *SegmentFileModelConfig
	LastSegmentIndex int
}

func CreateSegmentFileModel(config *SegmentFileModelConfig) *SegmentFileModel {
	segmentFileNameRegex := fmt.Sprintf(`%s_(\d*).seg`, config.FilePrefix)

	r := regexp.MustCompile(segmentFileNameRegex)
	dirEntries, err := os.ReadDir(config.Directory)
	if err != nil {
		panic(err)
	}
	var lastIndex = 0.0
	for _, entry := range dirEntries {

		match := r.FindStringSubmatch(entry.Name())
		// log.Debug().Msgf("reg %v v %v %v", segmentFileNameRegex, entry.Name(), len(match))
		if len(match) > 0 {
			segmentFileIndex, err := strconv.Atoi(match[1])
			if err != nil {
				panic(err)
			}
			lastIndex = math.Max(float64(segmentFileIndex), float64(lastIndex))
		}

	}

	return &SegmentFileModel{
		Config:           config,
		LastSegmentIndex: int(lastIndex),
	}
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

func readNextValue(f *os.File, sizeBufMaxBit int) ([]byte, error) {
	sizeBuf := make([]byte, sizeBufMaxBit)
	_, err := f.Read(sizeBuf)
	if err == io.EOF {
		return []byte{}, err
	}
	size := binary.LittleEndian.Uint64(sizeBuf)
	dataBuf := make([]byte, size)
	_, err = f.Read(dataBuf)
	if err == io.EOF {
		return []byte{}, err
	}
	return dataBuf[:], nil
}

func (s *SegmentFileModel) Flush(tuples []common.Tuple) error {
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

func (s *SegmentFileModel) Get(fileIndex int, pos int) (common.Tuple, error) {
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
	return common.Tuple{Key: keyByteArr, Value: valueArr}, nil
}
