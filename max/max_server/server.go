package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/ltbatista/findmaximum/max/maxpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) FindMaximum(stream maxpb.MaxService_FindMaximumServer) error {
	fmt.Printf("Findmaximum function was invoked with a streaming request\n")
	var max int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client streaming: %v", err)
			return err
		}
		number := req.GetMax().GetNumber()
		if number > max {
			max = number
		}
		result := strconv.Itoa(int(max))
		sendErr := stream.Send(&maxpb.MaxResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return sendErr
		}
	}
}

func main() {
	fmt.Println("Hello, I'm Max the Server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	maxpb.RegisterMaxServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
