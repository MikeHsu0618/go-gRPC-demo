package main

import (
	"log"
	"math/rand"
	"time"

	service2 "go-grpc-demo/calculator/client/service"
	pb "go-grpc-demo/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "127.0.0.1:50052"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v \n", err)
	}
	defer conn.Close()

	c := pb.NewCalculationServiceClient(conn)
	service := service2.NewService(c)

	// unary
	rand.Seed(time.Now().UnixNano())
	res := service.GetSum(int64(rand.Intn(30)), int64(rand.Intn(50)))
	log.Println(res.GetResult())

	// server stream
	service.GetPrime(int64(rand.Intn(50000)))

	// client stream
	service.GetAvg([]int64{1, 2, 3})

	// bidirectional stream
	service.GetMax([]int64{11, 23, 21, 34, 55, 77})

	// unary with error handle
	service.GetSqrt(-100)
}
