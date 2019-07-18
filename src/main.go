package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/musab/grpc-stream/src/protos"
)

type grpcStreamServer struct {
	savedNames []people // read-only after initialized
}

type people struct {
	FullName string   `json:"full_name"`
	Skill    []string `json:"skill"`
}

// loadPeople loads people from a JSON file.
func (s *grpcStreamServer) loadPeople() {
	var data []byte

	data = exampleData

	if err := json.Unmarshal(data, &s.savedNames); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

func (s *grpcStreamServer) GetNames(ctx context.Context, req *pb.GetNamesRequest) (*pb.GetNamesResponse, error) {
	names := make([]*pb.Name, len(s.savedNames))

	for i := range s.savedNames {
		names[i] = &pb.Name{FullName: s.savedNames[i].FullName}
	}

	return &pb.GetNamesResponse{Names: names}, nil
}

func (s *grpcStreamServer) ListSkills(req *pb.ListSkillsRequest, stream pb.GrpcStream_ListSkillsServer) error {
	for _, name := range s.savedNames {
		if err := stream.Send(&pb.ListSkillsResponse{Name: &pb.Name{FullName: name.FullName}, Skill: &pb.Skill{Language: name.Skill[0]}}); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

func newServer() *grpcStreamServer {
	s := &grpcStreamServer{}

	s.loadPeople()

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
	},
	{
    "skill": ["javascript"],
    "full_name": "Lisa Simpson"
  }
]
`)
