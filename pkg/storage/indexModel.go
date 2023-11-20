package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"nexus/pkg/common"
	"os"
	"regexp"
	"strconv"

	"github.com/rs/zerolog/log"
)

type IndexFileModelConfig struct {
	Directory  string
	FilePrefix string
}

// The index file contain fixed size entry, support fast binary search for the correct key position
type IndexFileModel struct {
	Config    *IndexFileModelConfig
	LastIndex int
}

func CreateIndexFileModel(config *IndexFileModelConfig) *IndexFileModel {
	indexFileNameRegex := fmt.Sprintf(`%s_(.*).index`, config.FilePrefix)
	r := regexp.MustCompile(indexFileNameRegex)
	dirEntries, err := os.ReadDir(config.Directory)
	if err != nil {
		panic(err)
	}
	var lastIndex = 0.0
	for _, entry := range dirEntries {
		match := r.FindStringSubmatch(entry.Name())
		if len(match) > 0 {
			indexFileIndex, err := strconv.Atoi(match[1])
			if err != nil {
				panic(err)
			}
			lastIndex = math.Max(float64(indexFileIndex), float64(lastIndex))
		}

	}

	return &IndexFileModel{
		Config:    config,
		LastIndex: int(lastIndex)}
}

func (i *IndexFileModel) Flush(tuples []common.Tuple) error {
	f, err := os.Create(i.getNextIndexPath())
	if err != nil {
		return err
	}
	defer f.Close()
	nextPos := 0
	for index := range tuples {
		_, err = f.Write([]byte(tuples[index].Key)[:KEY_SIZE_IN_BYTE])
		if err != nil {
			return err
		}
		nextPosByteArr := make([]byte, POSITION_OFFSET_SIZE_IN_BYTE)
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
		// TODO: Make this Binary search
		stat, err := f.Stat()
		if err != nil {
			panic(err)
		}
		fileSize := stat.Size()
		leftPos, rightPos := int64(0), fileSize/(KEY_SIZE_IN_BYTE+POSITION_OFFSET_SIZE_IN_BYTE)

		for {
			pos := (leftPos + rightPos) / 2 * (KEY_SIZE_IN_BYTE + POSITION_OFFSET_SIZE_IN_BYTE)
			_, err = f.Seek(pos, 0)
			if err != nil {
				panic(err)
			}
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
			if string(keyBuf) > string(key) {
				rightPos = (leftPos + rightPos) / 2
			} else {
				leftPos = (leftPos+rightPos)/2 + 1
			}
			// log.Info().Msgf("L %v R %v kb %v k %v  smaller %v i %v", leftPos, rightPos, string(keyBuf[:8]), string(key[:8]), string(keyBuf) < string(key), index)
			if leftPos >= rightPos {
				break
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
