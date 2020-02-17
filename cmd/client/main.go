package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dapr/go-sdk/dapr"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Client started")

	daprPort := os.Getenv("DAPR_GRPC_PORT")
	log.Printf("Connecting to Dapr on port [%v]", daprPort)
	daprAddress := fmt.Sprintf("localhost:%s", daprPort)
	conn, err := grpc.Dial(daprAddress, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Create the client
	client := dapr.NewDaprClient(conn)

	ctrl := newController(&client)

	r := mux.NewRouter()

	r.HandleFunc("/", ctrl.index)

	log.Println("Serving on port 8080...")
	http.ListenAndServe(":8080", r)
}

type controller struct {
	client dapr.DaprClient
}

func newController(c *dapr.DaprClient) *controller {
	return &controller{client: *c}
}

func (c *controller) index(w http.ResponseWriter, r *http.Request) {
	// Invoke a method called DoStuff on another Dapr enabled service with id client
	resp, err := c.client.InvokeService(context.Background(), &dapr.InvokeServiceEnvelope{
		Id:     "server",
		Data:   &any.Any{Value: []byte("Hello")},
		Method: "DoStuff",
	})

	if err != nil {
		log.Printf("Error calling server %v", err)
		http.Error(w, fmt.Sprintf("Error calling server %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, resp.GetData().String())
}
