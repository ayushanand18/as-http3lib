package main

import (
	"context"
	"log"

	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

type DummyResponse struct {
	Key   string `json:"key"`
	Value uint32 `json:"value"`
}

func HelloWorldGet(ctx context.Context, request interface{}) (response interface{}, err error) {
	return DummyResponse{
		Key:   "test",
		Value: 123,
	}, nil
}

func main() {
	ctx := context.Background()

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.GET("/json").Serve(types.ServeOptions{
		Handler: HelloWorldGet,
	})

	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
