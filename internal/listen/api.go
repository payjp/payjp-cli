package listen

import (
	"context"
	"errors"
	"io"
	"log"

	pb "github.com/payjp/payjp-cli/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartStream(ctx context.Context, address string, request *pb.ListenRequest, onEventHandler func(*pb.ListenResponse) error) error {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	grpcClient := pb.NewListenClient(conn)

	stream, err := grpcClient.Listen(ctx)
	if err != nil {
		log.Fatal("Failed to listen.")
		return err
	}

	err = stream.Send(request)
	if err != nil {
		log.Fatal("Failed to send request.")
		return err
	}

	for {
		received, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("server closed the stream")
				return err
			}

			log.Fatalf("Failed to receive event. %s", err)
			return err
		}

		err = onEventHandler(received)
		if err != nil {
			return err
		}
	}
}
