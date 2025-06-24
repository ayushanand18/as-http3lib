package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

type DummyResponse struct {
	Key   string `json:"key"`
	Value uint32 `json:"value"`
}

func main() {
	ctx := context.Background()

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.AddServeMethod(ctx, types.ServeOptions{
		URL:          "/json",
		ResponseType: constants.RESPONSE_TYPE_JSON_RESPONSE,
		Handler: func(ctx context.Context, r *http.Request) interface{} {
			return DummyResponse{
				Key:   "test",
				Value: 123,
			}
		},
		Method: "GET",
	})

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
