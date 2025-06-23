package main

import (
	"context"
	"io"
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
		Handler: func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "Hello World from GET.")
		},
		Method: "GET",
	})

	server.AddServeMethod(ctx, types.ServeOptions{
		URL: "/test",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "Hello World from POST.")
		},
		Method: "POST",
	})

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
