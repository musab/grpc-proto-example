package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/musab/grpc-stream/src/protos"
)

type grpcStreamServer struct {
	savedNames []*pb.Name // read-only after initialized
}

// loadFeatures loads features from a JSON file.
func (s *grpcStreamServer) loadFeatures() {
	var data []byte

	data = exampleData

	if err := json.Unmarshal(data, &s.savedNames); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

func (s *grpcStreamServer) GetNames(ctx context.Context, req *pb.GetNamesRequest) (*pb.GetNamesResponse, error) {
	names := make([]*pb.Name, len(s.savedNames))

	for i := range s.savedNames {
		names[i] = s.savedNames[i]
	}

	return &pb.GetNamesResponse{Names: names}, nil
}

func (s *grpcStreamServer) ListSkills(req *pb.ListSkillsRequest, stream pb.GrpcStream_ListSkillsServer) error {
	return errors.New("Not implemented")
}

func newServer() *grpcStreamServer {
	s := &grpcStreamServer{}

	s.loadFeatures()

	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGrpcStreamServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

var exampleData = []byte(`[
  {
    "skill": ["node"],
    "full_name": "Homer Simpson"
  },
  {
    "skill": ["c++"],
    "full_name": "Bart Simpson"
  },
  {
    "skill": ["go"],
    "full_name": "Marge Simpson"
  }
]
`)
