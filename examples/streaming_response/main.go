package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func HelloWorldStreaming(ctx context.Context, request interface{}) (response interface{}, err error) {
	for i := range 5 {
		time.Sleep(time.Duration(1) * time.Second)

		channel := ctx.Value(constants.STREAMING_RESPONSE_CHANNEL_CONTEXT_KEY).(chan types.StreamChunk)
		channel <- types.StreamChunk{
			Id:   uint32(i),
			Data: []byte(fmt.Sprintf("Chunk: %d \n\n", i)),
		}
	}

	return nil, nil
}

func main() {
	ctx := context.Background()

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.GET("/streaming").Serve(types.ServeOptions{
		Handler: HelloWorldStreaming,
		Options: types.MethodOptions{
			IsStreamingResponse: true,
		},
	})

	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
