package main

import (
	"fmt"
	"net"
	nexus "nexus/pkg/server"
	"os"

	pb "nexus/pkg/interface"
	"nexus/pkg/storage"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runServer() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	grpcServer := grpc.NewServer()
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
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
	log.Print("NExUS is listening on port 5555")
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func main() {
	runServer()
}
