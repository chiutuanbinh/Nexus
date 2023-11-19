package main

import (
	"nexus/pkg/storage"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func printAllEntry(s storage.Storage) {
	var lost = 0
	var cost int64 = 0
	for i := 1; i < 100000; i++ {
		key := strconv.Itoa(i)
		start := time.Now().UnixMilli()
		v, found := s.Get(key)
		cost += time.Now().UnixMilli() - start
		//we expect to find everything
		log.Info().Bool("found", found).Str("v", v).Str("K", key).Msg("")
		if !found {
			lost += 1
		}
		if i%1000 == 0 {
			log.Info().Msgf("HIT RATIO %v avg nano per query %v", float32(lost)/float32(i), cost/int64(i))
		}
	}
	log.Info().Int("lost", lost).Msg("")
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	s := storage.CreateStorage(&storage.StorageConfig{
		SSTableConfig: storage.SSTableConfig{
			Directory:        "/Users/bchiu/workspace/private/Nexus/data",
			FilePrefix:       "Nexus",
			SegmentThreshold: 1000000,
			MemtableMaxSize:  1000000,
			UseHash:          false,
		},
	})

	for i := 0; i < 100000; i++ {
		key := strconv.Itoa(i)
		x := time.Now().UnixMilli()
		err := s.Put(key, uuid.NewString())
		log.Info().Msgf("I time %v %v ", time.Now().UnixMilli()-x, i)
		if err != nil {
			panic(err)
		}
		_, found := s.Get(key)
		if !found {
			log.Error().Msgf("Not found %v", key)
		}
	}
	s.Flush()
	printAllEntry(s)
}
