package main

import (
	"context"
	"log"

	"github.com/ayushanand18/crazyhttp/pkg/types"
)

func main() {
	ctx := context.Background()

	server := crazyserver.NewHttpServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.GET("/test").Serve(types.ServeOptions{
		Handler: func(ctx context.Context, request interface{}) (interface{}, error) {
			return "Hello World from GET.", nil
		},
	})

	server.POST("/test").Serve(types.ServeOptions{
		Handler: func(ctx context.Context, request interface{}) (interface{}, error) {
			return "Hello World from POST.", nil
		},
	})

	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
