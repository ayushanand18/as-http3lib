package server

import (
	"github.com/ayushanand18/crazyhttp/pkg/types"
	web_socket "golang.org/x/net/websocket"
)

type websocket struct {
	Url string
	s   *server

	decoder               types.HttpEncoder
	encoder               types.HttpEncoder
	beforeServeMiddleware types.HttpRequestMiddleware
	afterServeMiddleware  types.HttpResponseMiddleware

	description string
	name        string
}

type WebSocket interface {
	Serve(types.HandlerFunc) WebSocket

	// Decoder for every message received
	WithDecoder(decoder types.HttpEncoder) WebSocket
	// Encoder for every message sent
	WithEncoder(encoder types.HttpEncoder) WebSocket
	// Middleware to run before every message is served
	WithBeforeServe(middleware types.HttpRequestMiddleware) WebSocket
	// Middleware to run after every message is sent
	WithAfterServe(middleware types.HttpResponseMiddleware) WebSocket
	// Name of the websocket endpoint - for Swagger API documentation
	WithName(name string) WebSocket
	// Description of the websocket endpoint - for Swagger API documentation
	WithDescription(desc string) WebSocket

	HandleHandshake(types.WebSocketHandshakeFunc) WebSocket
}

func NewWebsocket(url string, s *server) WebSocket {
	return &websocket{Url: url, s: s}
}

func (ws *websocket) Serve(handler types.HandlerFunc) WebSocket {
	fun := GetWebSocketHandlerFunc(handler)
	ws.s.mux.Handle(ws.Url, web_socket.Handler(fun))
	return ws
}

func (ws *websocket) WithDecoder(decoder types.HttpEncoder) WebSocket {
	ws.decoder = decoder
	return ws
}

func (ws *websocket) WithEncoder(encoder types.HttpEncoder) WebSocket {
	ws.encoder = encoder
	return ws
}

func (ws *websocket) WithBeforeServe(middleware types.HttpRequestMiddleware) WebSocket {
	ws.beforeServeMiddleware = middleware
	return ws
}

func (ws *websocket) WithAfterServe(middleware types.HttpResponseMiddleware) WebSocket {
	ws.afterServeMiddleware = middleware
	return ws
}

func (ws *websocket) WithName(name string) WebSocket {
	ws.name = name
	return ws
}

func (ws *websocket) WithDescription(desc string) WebSocket {
	ws.description = desc
	return ws
}

func (ws *websocket) HandleHandshake(fn types.WebSocketHandshakeFunc) WebSocket {
	return ws
}
