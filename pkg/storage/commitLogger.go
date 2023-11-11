package storage

import (
	"encoding/binary"
	"io"
	"os"
)

type CommitLogger interface {
	Write(key []byte, value []byte) error
	Load(readFunc func([]byte, []byte) error) error
	Clear() error
}

type SimpleCommitLogger struct {
	commitFileName string
}

func CreateCommitLog(commitFileName string) (CommitLogger, error) {
	return &SimpleCommitLogger{commitFileName: commitFileName}, nil
}

func (c *SimpleCommitLogger) Write(key []byte, value []byte) error {
	commitFile, err := os.OpenFile(c.commitFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return err
	}
	defer commitFile.Close()
	keySize := make([]byte, KEY_SIZE_MAX_BYTE_LENGTH)
	binary.LittleEndian.PutUint64(keySize, uint64(len(key)))
	_, err = commitFile.Write(keySize)
	if err != nil {
		return err
	}
	_, err = commitFile.Write(key)
	if err != nil {
		return err
	}
	valueSize := make([]byte, VALUE_SIZE_MAX_BYTE_LENGTH)
	binary.LittleEndian.PutUint64(valueSize, uint64(len(value)))
	_, err = commitFile.Write(valueSize)
	if err != nil {
		return err
	}
	_, err = commitFile.Write(value)
	if err != nil {
		return err
	}

	return nil
}

func (c *SimpleCommitLogger) Load(readFunc func([]byte, []byte) error) error {
	f, err := os.OpenFile(c.commitFileName, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	for {

		keySizeBuf := make([]byte, KEY_SIZE_MAX_BYTE_LENGTH)
		_, err := f.Read(keySizeBuf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		keySize := binary.LittleEndian.Uint64(keySizeBuf)
		keyBuf := make([]byte, keySize)
		_, err = f.Read(keyBuf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		valueSizeBuf := make([]byte, VALUE_SIZE_MAX_BYTE_LENGTH)
		_, err = f.Read(valueSizeBuf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		valueSize := binary.LittleEndian.Uint64(valueSizeBuf)
		valueBuf := make([]byte, valueSize)
		_, err = f.Read(valueBuf)

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = readFunc(keyBuf, valueBuf)
		if err != nil {
			return err
		}
	}
}

func (c *SimpleCommitLogger) Clear() error {
	err := os.Remove(c.commitFileName)
	return err
}
