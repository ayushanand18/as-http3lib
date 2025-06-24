package http3

import (
	"context"
	"net/http"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func streamingDefaultHandler(ctx context.Context, w http.ResponseWriter, options types.ServeOptions, handler types.HandlerFunc, r *http.Request) {

	for key, value := range options.DefaultHeaders {
		w.Header().Set(key, value)
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported by this server!", http.StatusInternalServerError)
		return
	}

	ch := make(chan types.StreamChunk)
	ctx = context.WithValue(ctx, constants.STREAMING_RESPONSE_CHANNEL_CONTEXT_KEY, ch)

	go func() {
		defer close(ch)
		handler(ctx, r)
	}()

	for chunk := range ch {
		if _, err := w.Write(chunk.Data); err != nil {
			break
		}

		flusher.Flush()
	}
}
