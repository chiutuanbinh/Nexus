package main

import (
	"fmt"
	"net"
	"nexus/pkg/config"
	nexus "nexus/pkg/server"
	"os"

	pb "nexus/pkg/interface"
	"nexus/pkg/storage"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runServer() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	viper.SetConfigName("nexus_config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error config file not found")
	}
	configLoader := func(c *config.Config) error {
		return viper.Unmarshal(c)
	}
	config.Init(configLoader)
	grpcServer := grpc.NewServer()
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(curDir + "/data"); os.IsNotExist(err) {
		err := os.Mkdir(curDir+"/data", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	storageImpl := storage.CreateStorage(&storage.StorageConfig{
		SSTableConfig: storage.SSTableConfig{
			Directory:        curDir + "/data",
			FilePrefix:       "Nexus",
			SegmentThreshold: 4 * 1024 * 1024,
			MemtableMaxSize:  1024 * 1024,
			UseHash:          true,
		}})
	server := nexus.CreateNexus(storageImpl)
	pb.RegisterNexusServer(grpcServer, server)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", "localhost:5555")
	if err != nil {
		fmt.Println(err)
		panic("yo port is taken")
	}
	log.Info().Msg("NExUS is listening on port 5555")
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func main() {
	runServer()
}
