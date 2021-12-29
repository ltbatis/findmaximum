package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ltbatista/findmaximum/max/maxpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm Max the Client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := maxpb.NewMaxServiceClient(cc)
	doBidiStreaming(c)
}

func doBidiStreaming(c maxpb.MaxServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream by invoking a client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := make([]*maxpb.MaxRequest, len(os.Args)-1)
	for i := 1; i < len(os.Args); i++ {
		numero, _ := strconv.Atoi(os.Args[i])
		requests[i-1] = &maxpb.MaxRequest{
			Max: &maxpb.Max{
				Number: int32(numero),
			},
		}
	}
	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving a response from server: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
