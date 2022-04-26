package main

import (
	"fmt"
	"log"
	"net"

	pb "go-grpc-demo/calculator/proto"
	"go-grpc-demo/calculator/server/service"
	"google.golang.org/grpc"
)

var addr = "127.0.0.1:50052"

func main() {
	lis, _ := net.Listen("tcp", addr)
	fmt.Println("listening on ", addr)

	s := grpc.NewServer()
	pb.RegisterCalculationServiceServer(s, service.NewService())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v \n", err)
	}

}
