package internal

import (
	"context"

	pb "nexus/proto-gen/nexus"

	"google.golang.org/protobuf/types/known/emptypb"
)

type NexusServer struct {
	pb.UnimplementedNexusServer
	storage map[string]string
}

func CreateNexus() *NexusServer {
	return &NexusServer{storage: map[string]string{}}
}

func (s *NexusServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *NexusServer) Put(ctx context.Context, putRequest *pb.PutRequest) (*pb.PutReponse, error) {
	s.storage[putRequest.Key] = putRequest.Value
	return &pb.PutReponse{Err: 0}, nil
}

func (s *NexusServer) Get(ctx context.Context, getRequest *pb.GetRequest) (*pb.GetReponse, error) {
	return &pb.GetReponse{Err: 0, Key: getRequest.Key, Value: s.storage[getRequest.Key]}, nil
}
