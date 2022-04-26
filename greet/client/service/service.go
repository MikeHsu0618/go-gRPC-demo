package service

import (
	"context"
	"io"
	"log"
	"time"

	pb "go-grpc-demo/greet/proto"
)

var ctx = context.Background()

type Service interface {
	GetGreet(name string) *pb.GreetResponse
	GetGreetManyTimes(name string)
	GetLongGreet(names []string)
	GetGreetEveryone(names []string)
}

type service struct {
	pb.GreetServiceClient
}

func NewService(c pb.GreetServiceClient) Service {
	return &service{c}
}

func (s *service) GetGreet(name string) *pb.GreetResponse {
	r, _ := s.Greet(ctx, &pb.GreetRequest{FirstName: name})
	return r
}

func (s *service) GetGreetManyTimes(name string) {
	req := &pb.GreetRequest{FirstName: name}
	stream, _ := s.GreetManyTimes(ctx, req)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading the stream: %v \n", err)
		}

		log.Printf("GreetManyTimes : %s \n", msg.GetResult())
	}
}

func (s *service) GetLongGreet(names []string) {
	log.Println("GetLongGreet was invoked")

	// reqs := []*pb.GreetRequest{
	// 	{FirstName: "test1"},
	// 	{FirstName: "test2"},
	// 	{FirstName: "test3"},
	// }
	var reqs []*pb.GreetRequest
	for _, name := range names {
		reqs = append(reqs, &pb.GreetRequest{FirstName: name})
	}

	stream, err := s.LongGreet(ctx)
	if err != nil {
		log.Fatalf("Error while reading the stream: %v \n", err)
	}

	for _, req := range reqs {
		log.Printf("Sending req: %v\n", req)
		stream.Send(req)

		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving from LongGreet: %v\n", err)
	}

	log.Printf("LongGreet: %#v\n", res.GetResult())
}

func (s *service) GetGreetEveryone(names []string) {
	log.Println("GetGreetEveryone was invoked")

	stream, err := s.GreetEveryone(ctx)

	if err != nil {
		log.Fatalf("Error while creaeting stream: %v \n", err)
	}

	var reqs []*pb.GreetRequest
	for _, name := range names {
		reqs = append(reqs, &pb.GreetRequest{FirstName: name})
	}

	waitChan := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Printf("Send request : %v \n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("Eroor while receiving: %v \n ", err)
			}

			log.Printf("Received: %v \n ", res.GetResult())
		}

		close(waitChan)
	}()

	<-waitChan
}
