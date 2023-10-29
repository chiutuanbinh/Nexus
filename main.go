package main

import (
	"fmt"
	"net"
	nexus "nexus/internal"
	"nexus/internal/crawler"

	pb "nexus/internal/interface"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runServer() {
	grpcServer := grpc.NewServer()
	server := nexus.CreateNexus()
	pb.RegisterNexusServer(grpcServer, server)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", "localhost:5555")
	if err != nil {
		fmt.Println(err)
		panic("yo port is taken")
	}
	grpcServer.Serve(lis)
}

func main() {
	crawler.Foo()
}
