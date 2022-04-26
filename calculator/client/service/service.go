package service

import (
	"context"
	"io"
	"log"
	"time"

	pb "go-grpc-demo/calculator/proto"
)

var ctx = context.Background()

type Service interface {
	GetSum(firstNum int64, secondNum int64) *pb.SumResponse
	GetPrime(num int64)
	GetAvg(nums []int64)
	GetMax(nums []int64)
}

type service struct {
	pb.CalculationServiceClient
}

func NewService(c pb.CalculationServiceClient) Service {
	return &service{c}
}

func (s *service) GetSum(firstNum int64, secondNum int64) *pb.SumResponse {
	log.Println("get Sum !", firstNum, secondNum)
	r, _ := s.Sum(ctx, &pb.SumRequest{
		FirstNumber:  firstNum,
		SecondNumber: secondNum,
	})
	return r
}

func (s *service) GetPrime(num int64) {
	req := &pb.PrimeRequest{Number: num}
	stream, _ := s.Prime(ctx, req)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading the stream: %v \n", err)
		}

		log.Printf("GetPrime : %v \n", msg.GetResult())
		time.Sleep(time.Second)
	}
}

func (s *service) GetAvg(nums []int64) {
	log.Println("GetAvg was invoked")

	stream, err := s.Avg(ctx)
	if err == io.EOF {
		log.Fatalf("Error while reading the stream: %v \n", err)
	}

	for _, num := range nums {
		log.Printf("Sending req: %v", num)
		stream.Send(&pb.AvgRequest{Number: num})
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving from GetAvg: %v\n", err)
	}

	log.Printf("GetAvg: %#v\n", res.GetResult())
}

func (s *service) GetMax(nums []int64) {
	log.Println("GetAvg was invoked")

	stream, err := s.Max(ctx)

	if err == io.EOF {
		log.Fatalf("Error while reading the stream: %v \n", err)
	}

	waitChan := make(chan struct{})
	go func() {
		for _, num := range nums {
			stream.Send(&pb.MaxRequest{Number: num})
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

			log.Printf("Current Max Num Received: %v \n ", res.GetResult())
		}
		close(waitChan)
	}()

	<-waitChan
}
