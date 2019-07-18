package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/musab/grpc-stream/src/protos"
	"google.golang.org/grpc"
)

func printName(client pb.GrpcStreamClient) {
	log.Printf("Getting names")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	names, err := client.GetNames(ctx, &pb.GetNamesRequest{})
	if err != nil {
		log.Fatalf("%v.GetNames(_) = _, %v: ", client, err)
	}
	log.Println(names)
}

func listSkills(client pb.GrpcStreamClient) {
	stream, err := client.ListSkills(context.Background(), &pb.ListSkillsRequest{})

	if err != nil {
		log.Fatalf("%v.listSkills(_) = _, %v", client, err)
	}

	for {
		name, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.listSkills(_) = _, %v", client, err)
		}
		log.Println(name)
	}

}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:7777", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewGrpcStreamClient(conn)

	printName(client)
	listSkills(client)
}
