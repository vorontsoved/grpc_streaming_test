package main

import (
	"bufio"
	"context"
	"google.golang.org/grpc"
	pb "grpc_streaming_test/generated/code"
	"io"
	"log"
	"os"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamingServiceClient(conn)

	stream, err := client.BidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
			}
			log.Printf("Received response from server: %s", response.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if err := stream.Send(&pb.StreamRequest{Message: message}); err != nil {
			log.Fatalf("Error while sending: %v", err)
		}
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Error while closing stream: %v", err)
	}

	<-waitc
}
