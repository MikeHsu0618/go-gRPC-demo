package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"

	pb "go-grpc-demo/calculator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	pb.UnimplementedCalculationServiceServer
}

func NewService() pb.CalculationServiceServer {
	return &service{pb.UnimplementedCalculationServiceServer{}}
}

func (s *service) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	log.Printf("Get First Number : %v, Second Number : %v", in.GetFirstNumber(), in.SecondNumber)
	return &pb.SumResponse{Result: in.GetFirstNumber() + in.GetSecondNumber()}, nil
}

func (s *service) Sqrt(ctx context.Context, in *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	num := in.GetNumber()
	log.Printf("Get Number : %v", num)

	if num < 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %d", num),
		)
	}
	return &pb.SqrtResponse{Result: math.Sqrt(float64(num))}, nil
}

func (s *service) Prime(in *pb.PrimeRequest, stream pb.CalculationService_PrimeServer) error {
	log.Printf("primes function was invoked with %v", in.GetNumber())

	num := in.GetNumber()
	divisor := int64(2)
	for num > 1 {
		if num%divisor == 0 {
			stream.Send(&pb.PrimeResponse{Result: divisor})
			num /= divisor
		} else {
			divisor++
		}
	}

	return nil
}

func (s *service) Avg(stream pb.CalculationService_AvgServer) error {
	log.Println("Avg function was invoked")
	var sum int64
	var count int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := float64(sum) / float64(count)
			return stream.SendAndClose(&pb.AvgResponse{Result: res})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream")
		}

		log.Println("Receiving number: ", req.GetNumber())
		count++
		sum += req.GetNumber()
	}
}

func (s *service) Max(stream pb.CalculationService_MaxServer) error {
	log.Println("Max function was invoked")
	max := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream")
		}

		num := req.GetNumber()
		if num > max {
			max = num
		}
		err = stream.Send(&pb.MaxResponse{Result: max})

		if err != nil {
			log.Fatalf("Error while sending data to client : %v\n", err)
		}
	}
}
