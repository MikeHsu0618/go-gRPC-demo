package main

import (
	"log"
	"net"

	pb "go-grpc-demo/greet/proto"
	"go-grpc-demo/greet/server/service"
	"google.golang.org/grpc"
)

var addr = "127.0.0.1:50051"

func main() {
	lis, _ := net.Listen("tcp", addr)

	log.Println("listening on", addr)

	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, service.NewService())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v \n", err)
	}
}
