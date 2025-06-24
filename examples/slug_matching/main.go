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
		URL: "/users/{user_id}",
		Handler: func(ctx context.Context, r *http.Request) *types.HttpResponse {
			headers := make(map[string]string)
			headers["X-User-Id"] = r.PathValue("user_id")

			return &types.HttpResponse{
				StatusCode: 200,
				Headers:    headers,
				Body:       []byte("Hello World from GET."),
			}
		},
		Method: "GET",
	})

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
