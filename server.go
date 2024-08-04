package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "grpc_streaming_test/generated/code"
)

type server struct {
	pb.UnimplementedStreamingServiceServer
}

func (s *server) BidirectionalStreaming(stream pb.StreamingService_BidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Received message from client: %s", req.Message)
		response := &pb.StreamResponse{
			Message: "Received: " + req.Message,
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamingServiceServer(s, &server{})

	log.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
