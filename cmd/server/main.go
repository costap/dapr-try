package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dapr/go-sdk/daprclient"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type server struct {
}

func main() {
	log.Println("Server starting...")

	// create listiner
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	s := grpc.NewServer()
	daprclient.RegisterDaprClientServer(s, &server{})

	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) DoStuff() string {
	return "Hello from server!"
}

// This method gets invoked when a remote service has called the app through Dapr
// The payload carries a Method to identify the method, a set of metadata properties and an optional payload
func (s *server) OnInvoke(ctx context.Context, in *daprclient.InvokeEnvelope) (*any.Any, error) {
	var response string

	log.Println(fmt.Sprintf("Got invoked with: %s", string(in.Data.Value)))

	switch in.Method {
	case "DoStuff":
		response = s.DoStuff()
	}
	return &any.Any{
		Value: []byte(response),
	}, nil
}

// Dapr will call this method to get the list of topics the app wants to subscribe to. In this example, we are telling Dapr
// To subscribe to a topic named TopicA
func (s *server) GetTopicSubscriptions(ctx context.Context, in *empty.Empty) (*daprclient.GetTopicSubscriptionsEnvelope, error) {
	log.Println("GetTopicSubscriptions called.")
	return &daprclient.GetTopicSubscriptionsEnvelope{
		Topics: []string{"TopicA"},
	}, nil
}

// Dapper will call this method to get the list of bindings the app will get invoked by. In this example, we are telling Dapr
// To invoke our app with a binding named storage
func (s *server) GetBindingsSubscriptions(ctx context.Context, in *empty.Empty) (*daprclient.GetBindingsSubscriptionsEnvelope, error) {
	log.Println("GetBindingsSubscriptions called.")
	return &daprclient.GetBindingsSubscriptionsEnvelope{
		Bindings: []string{"storage"},
	}, nil
}

// This method gets invoked every time a new event is fired from a registerd binding. The message carries the binding name, a payload and optional metadata
func (s *server) OnBindingEvent(ctx context.Context, in *daprclient.BindingEventEnvelope) (*daprclient.BindingResponseEnvelope, error) {
	log.Println("Invoked from binding")
	return &daprclient.BindingResponseEnvelope{}, nil
}

// This method is fired whenever a message has been published to a topic that has been subscribed. Dapr sends published messages in a CloudEvents 0.3 envelope.
func (s *server) OnTopicEvent(ctx context.Context, in *daprclient.CloudEventEnvelope) (*empty.Empty, error) {
	log.Println("Topic message arrived")
	return &empty.Empty{}, nil
}