package listen

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/payjp/payjp-cli/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Listener struct {
	address   string
	client    pb.ListenClient
	connected bool
}

var ReconnectRequiredError = errors.New("reconnect required")

func NewListener(address string) *Listener {
	return &Listener{
		address:   address,
		connected: false,
	}
}

func (l *Listener) StartListen(ctx context.Context, request *pb.ListenRequest, onEventHandler func(*pb.PayjpEventResponse) error) error {
	conn, err := grpc.NewClient(
		l.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	l.client = pb.NewListenClient(conn)
	log.Println("Connecting to server...")

	for {
		err = l.listen(ctx, request, onEventHandler)
		if err != nil {
			if errors.Is(err, ReconnectRequiredError) {
				log.Println("Reconnecting...")
				continue
			}

			return err
		}
	}
}

func (l *Listener) listen(ctx context.Context, request *pb.ListenRequest, onEventHandler func(*pb.PayjpEventResponse) error) error {
	receiveCh := make(chan *pb.ListenResponse)
	errCh := make(chan error)
	timeoutCh := make(chan struct{})

	stream, err := l.client.Listen(ctx)
	if err != nil {
		log.Fatal("Failed to listen.")
		return err
	}

	err = stream.Send(request)
	if err != nil {
		log.Fatal("Failed to send request.")
		return err
	}

	go func() {
		time.Sleep(5 * time.Second)
		timeoutCh <- struct{}{}
	}()

	go func() {
		for {
			received, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			receiveCh <- received
		}
	}()

	for {
		select {
		case <-timeoutCh:
			if l.connected {
				continue
			}
			log.Println("Connection timeout, Please try again later.")
			return fmt.Errorf("timeout")
		case err := <-errCh:
			if errors.Is(err, io.EOF) {
				log.Println("Server closed the stream")
				l.connected = false
				return ReconnectRequiredError
			}

			log.Fatalf("Failed to receive event. %s", err)
			return err
		case received := <-receiveCh:
			switch received.Response.(type) {
			case *pb.ListenResponse_PayjpEventResponse:
				err = onEventHandler(received.GetPayjpEventResponse())
				if err != nil {
					return err
				}
			case *pb.ListenResponse_SystemEventResponse:
				switch received.GetSystemEventResponse().Type {
				case pb.SystemEventType_SYSTEM_EVENT_TYPE_OK:
					l.connected = true
					log.Println("Connected. Start listening...")
				case pb.SystemEventType_SYSTEM_EVENT_TYPE_RECONNECT_REQUESTED:
					return ReconnectRequiredError
				}
			}
		}
	}
}
