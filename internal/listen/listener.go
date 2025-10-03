package listen

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	pb "github.com/payjp/payjp-cli/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const MaxReconnectAttempts = 3

type Listener struct {
	address           string
	client            pb.ListenClient
	connected         bool
	reconnectAttempts int
}

var ReconnectRequiredError = errors.New("reconnect required")

func NewListener(address string) *Listener {
	return &Listener{
		address:           address,
		connected:         false,
		reconnectAttempts: 0,
	}
}

func (l *Listener) StartListen(ctx context.Context, request *pb.ListenRequest, onEventHandler func(*pb.PayjpEventResponse) error) error {
	cred := credentials.NewTLS(&tls.Config{})
	if !strings.HasSuffix(l.address, ":443") {
		cred = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(
		l.address,
		grpc.WithTransportCredentials(cred),
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
				l.reconnectAttempts++
				if l.reconnectAttempts <= MaxReconnectAttempts {
					log.Printf("Reconnecting... (%d/%d)", l.reconnectAttempts, MaxReconnectAttempts)
					continue
				} else {
					return fmt.Errorf("max reconnect attempts reached")
				}
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
			log.Println("Connection timeout, retrying...")
			return ReconnectRequiredError
		case err := <-errCh:
			l.connected = false
			if stat, ok := status.FromError(err); ok {
				if stat.Code() == codes.Unauthenticated {
					return fmt.Errorf("authentication failed. Please login again and try your request.")
				} else if stat.Code() == codes.FailedPrecondition {
					return fmt.Errorf("%s", stat.Message())
				} else if stat.Code() == codes.Internal || stat.Code() == codes.Unavailable {
					return ReconnectRequiredError
				}
			}
			if errors.Is(err, io.EOF) {
				log.Println("Server closed the stream")
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
					l.reconnectAttempts = 0
					log.Println("Connected. Start listening...")
				case pb.SystemEventType_SYSTEM_EVENT_TYPE_RECONNECT_REQUESTED:
					return ReconnectRequiredError
				case pb.SystemEventType_SYSTEM_EVENT_TYPE_PING:
					stream.Send(&pb.ListenRequest{
						Request: &pb.ListenRequest_PongRequest{
							PongRequest: &pb.PongRequest{
								Timestamp: received.GetSystemEventResponse().PingData.Timestamp,
							},
						},
					})
				}
			}
		}
	}
}
