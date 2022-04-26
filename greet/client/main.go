package main

import (
	"log"

	service2 "go-grpc-demo/greet/client/service"
	pb "go-grpc-demo/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "localhost:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v \n", err)
	}
	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)
	service := service2.NewService(c)

	// unary
	res := service.GetGreet("foo")
	log.Println(res.GetResult())

	// server stream
	service.GetGreetManyTimes("bar")

	// client stream
	service.GetLongGreet([]string{"1", "2", "3"})

	// bidirectional stream
	service.GetGreetEveryone([]string{"name1", "name2", "name3"})
}
