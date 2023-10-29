package main

import (
	"fmt"
	"log"
	"net"
	nexus "nexus/internal"
	"os"

	pb "nexus/internal/interface"
	"nexus/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runServer() {
	grpcServer := grpc.NewServer()
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
	server := nexus.CreateNexus(sstable, commitLogger)
	pb.RegisterNexusServer(grpcServer, server)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", "localhost:5555")
	if err != nil {
		fmt.Println(err)
		panic("yo port is taken")
	}
	log.Println("NExUS is listening on port 5555")
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func main() {
	runServer()
}
