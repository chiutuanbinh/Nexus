package server

import (
	"context"

	pb "nexus/pkg/interface"
	"nexus/pkg/storage"

	"google.golang.org/protobuf/types/known/emptypb"
)

type NexusGrpcServer struct {
	pb.UnimplementedNexusServer
	storage storage.Storage
}

func CreateNexus(storage storage.Storage) *NexusGrpcServer {
	return &NexusGrpcServer{storage: storage}
}

func (s *NexusGrpcServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *NexusGrpcServer) Put(ctx context.Context, putRequest *pb.PutRequest) (*pb.PutResponse, error) {
	if s.storage.Put(putRequest.GetKey(), putRequest.GetValue()) != nil {
		return &pb.PutResponse{Err: -1}, nil
	}
	return &pb.PutResponse{Err: 0}, nil
}

func (s *NexusGrpcServer) Get(ctx context.Context, getRequest *pb.GetRequest) (*pb.GetResponse, error) {
	value, found := s.storage.Get(getRequest.GetKey())
	if !found {
		return &pb.GetResponse{Err: -1}, nil
	}
	return &pb.GetResponse{Err: 0, Key: getRequest.Key, Value: value}, nil
}

func (s *NexusGrpcServer) Delete(ctx context.Context, deleteRequest *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.storage.Delete(deleteRequest.GetKey())
	if err != nil {
		return &pb.DeleteResponse{Err: -1}, nil
	}
	return &pb.DeleteResponse{Err: 0}, nil
}

func (s *NexusGrpcServer) Flush(ctx context.Context, flushRequest *pb.FlushRequest) (*pb.FlushResponse, error) {
	err := s.storage.Flush()
	if err != nil {
		return &pb.FlushResponse{Err: -1}, nil
	}
	return &pb.FlushResponse{Err: 0}, nil
}
