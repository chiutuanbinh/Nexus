package storage

import (
	storage "nexus/pkg/storage/internal"
	"os"
)

type Storage interface {
	Put(key string, value string) error
	Get(key string) (string, bool)
	Delete(key string) error
	Flush() error
}
type storageImplementation struct {
	commitLogger storage.CommitLogger
	sstable      *storage.SSTable
}

// Delete implements Storage.
func (s *storageImplementation) Delete(key string) error {
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

func CreateStorage() Storage {
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataDirectory := curDir + "/data"
	commitLogFilePath := dataDirectory + "/commit.log"
	commitLogger, err := storage.CreateCommitLog(commitLogFilePath)
	if err != nil {
		panic(err)
	}

	sstable := storage.NewSSTable(&storage.SSTableConfig{
		Directory:        curDir + "/data",
		FilePrefix:       "Nexus",
		SegmentThreshold: 4 * 1024 * 1024,
		MemtableMaxSize:  1024 * 1024,
		UseHash:          true,
	}, commitLogger.Clear)
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
