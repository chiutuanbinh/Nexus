package main

import (
	"nexus/pkg/storage"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	s := storage.CreateStorage(&storage.StorageConfig{
		SSTableConfig: storage.SSTableConfig{
			Directory:        "data",
			FilePrefix:       "Nexus",
			SegmentThreshold: 10000,
			MemtableMaxSize:  10000,
			UseHash:          true,
		},
	})

	for i := 0; i < 100000; i++ {
		key := strconv.Itoa(i)
		err := s.Put(key, uuid.NewString())
		if err != nil {
			panic(err)
		}
	}
	s.Flush()

	for i := 0; i < 100000; i++ {
		key := strconv.Itoa(i)
		v, found := s.Get(key)
		log.Info().Bool("found", found).Str("v", v).Str("K", key).Msg("")
	}

}
