package http3

import (
	"context"
	"net/http"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	ashttp "github.com/ayushanand18/as-http3lib/internal/http"
	"github.com/ayushanand18/as-http3lib/pkg/errors"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func streamingDefaultHandler(
	ctx context.Context,
	w http.ResponseWriter,
	handler types.HandlerFunc,
	decoder types.HttpDecoder,
	encoder types.HttpEncoder,
	r *http.Request) {

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
		var request interface{}
		var err error

		if decoder != nil {
			request, err = decoder(ctx, r)
			if err != nil {
				w.WriteHeader(errors.DecodeErrorToHttpErrorStatus(err))
				return
			}
		}

		_, err = handler(ctx, request)
		if err != nil {
			w.WriteHeader(errors.DecodeErrorToHttpErrorStatus(err))
			return
		}

	}()

	for chunk := range ch {
		headers := make(map[string][]string)
		var encoded []byte
		var err error

		if encoder != nil {
			headers, encoded, err = encoder(ctx, chunk.Data)
			if err != nil {
				w.WriteHeader(errors.DecodeErrorToHttpErrorStatus(err))
				break
			}
		} else {
			headers, encoded, err = ashttp.DefaultHttpEncode(ctx, chunk.Data)
			if err != nil {
				w.WriteHeader(errors.DecodeErrorToHttpErrorStatus(err))
				break
			}
		}

		for key, value := range headers {
			w.Header().Del(key)
			for _, v := range value {
				w.Header().Add(key, v)
			}
		}

		if _, err := w.Write(encoded); err != nil {
			break
		}

		flusher.Flush()
	}
}
