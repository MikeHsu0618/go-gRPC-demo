package service

import (
	"context"
	"io"
	"log"

	pb "go-grpc-demo/calculator/proto"
)

// type Service interface {
// 	Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error)
// 	Prime(in *pb.PrimeRequest, stream pb.CalculationService_PrimeServer) error
// 	Avg(stream pb.CalculationService_AvgServer) error
// }

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
