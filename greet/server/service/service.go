package service

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "go-grpc-demo/greet/proto"
)

// type Server interface {
// 	Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error)
// 	GreetManyTimes(in *pb.GreetRequest, stream pb.GreetService_GreetManyTimesServer) error
// 	LongGreet(stream pb.GreetService_LongGreetServer) error
// 	GreetEveryone(pb.GreetService_GreetEveryoneServer) error
// }

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreetServiceServer
}

func NewService() pb.GreetServiceServer {
	return &server{pb.UnimplementedGreetServiceServer{}}
}

// SayHello implements helloworld.GreeterServer
func (s *server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Received: %v", in.GetFirstName())
	return &pb.GreetResponse{Result: "Hello " + in.GetFirstName()}, nil
}

func (s *server) GreetManyTimes(in *pb.GreetRequest, stream pb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes function was invoked with: %v \n", in)
	for i := 0; i < 10; i++ {
		res := fmt.Sprintf("Hello %s, number %d", in.FirstName, i)
		err := stream.Send(&pb.GreetResponse{Result: res})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	log.Println("Long Greet function was invoked")

	res := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.GreetResponse{Result: res})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream %v\n", err)
		}

		res += fmt.Sprintf("Hello %s \n", req.GetFirstName())
	}
}

func (s *server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	log.Println("GreetEveryone function was invoked")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream %v\n", err)
		}

		res := "Hello " + req.GetFirstName() + " !"
		err = stream.Send(&pb.GreetResponse{Result: res})

		if err != nil {
			log.Fatalf("Error while sending data to client : %v\n", err)
		}
	}
}
