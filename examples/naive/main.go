package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func main() {
	ctx := context.Background()

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.AddServeMethod(ctx, types.ServeOptions{
		URL: "/test",
		Handler: func(ctx context.Context, r *http.Request) (interface{}, error) {
			return "Hello World from GET.", nil
		},
		Method: "GET",
	})

	server.AddServeMethod(ctx, types.ServeOptions{
		URL: "/test",
		Handler: func(ctx context.Context, r *http.Request) (interface{}, error) {
			return "Hello World from POST.", nil
		},
		Method: "POST",
	})

	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
