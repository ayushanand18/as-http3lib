package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ayushanand18/as-http3lib/internal/constants"
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
		URL:          "/streaming",
		ResponseType: constants.RESPONSE_TYPE_STREAMING_RESPONSE,
		Handler: func(ctx context.Context, r *http.Request) interface{} {
			for i := range 5 {
				time.Sleep(time.Duration(1) * time.Second)

				ctx.Value(constants.STREAMING_RESPONSE_CHANNEL_CONTEXT_KEY).(chan types.StreamChunk) <- types.StreamChunk{
					Id:   uint32(i),
					Data: []byte(fmt.Sprintf("Chunk: %d \n\n", i)),
				}
			}

			return nil
		},
		Method: "GET",
	})

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
