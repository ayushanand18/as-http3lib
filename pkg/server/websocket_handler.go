package server

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ayushanand18/crazyhttp/internal/constants"
	ashttp "github.com/ayushanand18/crazyhttp/internal/http"
	"github.com/ayushanand18/crazyhttp/pkg/errors"
	"github.com/ayushanand18/crazyhttp/pkg/types"
	gws "github.com/gorilla/websocket"
)

func websocketHandler(
	ctx context.Context,
	conn *gws.Conn,
	w http.ResponseWriter,
	r *http.Request,
	ws *websocket,
	handler types.WebsocketHandlerFunc,
) {

	requestChannel := make(chan types.WebsocketStreamChunk)
	ctx = context.WithValue(ctx, constants.WebsocketRequestChannel, requestChannel)

	responseChannel := make(chan types.WebsocketStreamChunk)
	ctx = context.WithValue(ctx, constants.WebsocketResponseChannel, responseChannel)

	go func() {
		defer close(requestChannel)
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				if err == io.EOF {
					slog.Info("WebSocket connection closed by client")
				} else {
					slog.Error("Error receiving WebSocket message", "error", err)
				}
				break
			}

			if ws.rateLimiter != nil {
				key := ctx.Value(constants.RateLimitCustomKey)
				if key == nil || key == "" {
					key = strings.Split(r.RemoteAddr, ":")[0]
				}
				_, ok := key.(string)
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
					slog.ErrorContext(ctx, "rate limit key is not a string", "key:=", key)
					return
				}
				ws.rateLimiter.Allow(key.(string))
			}

			msg, err := ashttp.GetDefaultSerialization(message)
			if err != nil {
				w.WriteHeader(errors.DecodeErrorToHttpErrorStatus(err))
				return
			}

			requestChannel <- types.WebsocketStreamChunk{MessageType: mt, Data: msg}
		}
	}()

	go func() {
		handler(ctx)
	}()

	for chunk := range responseChannel {
		if chunk.MessageType == 0 {
			chunk.MessageType = 1
		}

		var headers map[string][]string
		var encoded []byte
		var err error

		if ws.encoder != nil {
			headers, encoded, err = ws.encoder(ctx, chunk.Data)
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

		err = conn.WriteMessage(chunk.MessageType, encoded)
		if err != nil {
			slog.Error("Error sending WebSocket message", "error", err)
			break
		}

	}
}
