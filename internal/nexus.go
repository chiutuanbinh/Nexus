package internal

import (
	"context"

	pb "nexus/internal/interface"
	"nexus/internal/storage"

	"google.golang.org/protobuf/types/known/emptypb"
)

type NexusServer struct {
	pb.UnimplementedNexusServer
	storage *storage.SSTable
}

func CreateNexus() *NexusServer {
	return &NexusServer{storage: storage.NewSSTable(&storage.SSTableConfig{
		Directory:        "data",
		FilePrefix:       "Nexus",
		SegmentThreshold: 4 * 1024 * 1024,
		MemtableMaxSize:  1024 * 1024,
		UseHash:          true,
	})}
}

func (s *NexusServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *NexusServer) Put(ctx context.Context, putRequest *pb.PutRequest) (*pb.PutResponse, error) {
	err := s.storage.Insert(putRequest.Key, putRequest.Value)
	if err != nil {
		return nil, err
	}
	return &pb.PutResponse{Err: 0}, nil
}

func (s *NexusServer) Get(ctx context.Context, getRequest *pb.GetRequest) (*pb.GetResponse, error) {
	value, found := s.storage.Find(getRequest.GetKey())
	if !found {
		return &pb.GetResponse{Err: -1}, nil
	}
	return &pb.GetResponse{Err: 0, Key: getRequest.Key, Value: value}, nil
}
