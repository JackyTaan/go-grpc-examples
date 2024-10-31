package main

import (
	"io"
	"log"
	"net"

	"github.com/jackytaan/go-grpc-examples/stream/bi-directional-streaming/feeds/feedpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	} else {
		log.Println("Server is listening on", "localhost:50051")
	}

	s := grpc.NewServer()
	feedpb.RegisterFeedsServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}

// Broadcast reads client stream and broadcasts recieved feeds
func (*server) Broadcast(stream feedpb.Feeds_BroadcastServer) error {
	iCount := 0
	for {
		msg, err := stream.Recv()
		iCount++
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("could not recieve from stream : %v", err)
			return err
		}
		feed := msg.GetFeed()
		// fmt.Println("Server recieved: ", feed)
		//feed = feed + "-" + msg.GetFeed()
		//fmt.Println("Server sent to client: ", feed)
		stream.Send(&feedpb.FeedResponse{Feed: feed})

	}

}
