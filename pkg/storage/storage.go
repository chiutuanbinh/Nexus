package storage

import (
	"os"
)

type Storage interface {
	Put(key string, value string) error
	Get(key string) (string, bool)
	Delete(key string) error
	Flush() error
}
type storageImplementation struct {
	commitLogger CommitLogger
	sstable      *SSTable
}

type StorageConfig struct {
	SSTableConfig
}

// Delete implements Storage.
func (s *storageImplementation) Delete(key string) error {
	err := s.commitLogger.Write([]byte(key), []byte(TOMBSTONE))
	if err != nil {
		return err
	}
	return s.sstable.Delete(key)
}

// Flush implements Storage.
func (s *storageImplementation) Flush() error {
	return s.sstable.Flush()
}

// Get implements Storage.
func (s *storageImplementation) Get(key string) (string, bool) {
	return s.sstable.Find(key)
}

// Put implements Storage.
func (s *storageImplementation) Put(key string, value string) error {
	err := s.commitLogger.Write([]byte(key), []byte(value))
	if err != nil {
		return err
	}
	err = s.sstable.Insert(key, value)
	return err
}

func CreateStorage(config *StorageConfig) Storage {
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataDirectory := curDir + "/data"
	commitLogFilePath := dataDirectory + "/commit.log"
	commitLogger, err := CreateCommitLog(commitLogFilePath)
	if err != nil {
		panic(err)
	}

	sstable := NewSSTable(&config.SSTableConfig, commitLogger.Clear)
	commitlogConsumer := func(key []byte, value []byte) error {
		return sstable.Insert(string(key), string(value))
	}
	err = commitLogger.Load(commitlogConsumer)
	if err != nil {
		panic(err)
	}
	return &storageImplementation{
		commitLogger: commitLogger,
		sstable:      sstable,
	}
}
