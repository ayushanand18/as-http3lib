package main_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
	qchttp3 "github.com/quic-go/quic-go/http3"
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

func TestHTTP3Server_BasicStreamingResponse(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:4431"

	s := http3.NewServer(ctx)
	if err := s.Initialize(ctx); err != nil {
		t.Fatalf("server initialization failed: %v", err)
	}

	s.GET("/streaming").Serve(types.ServeOptions{
		Handler: HelloWorldStreaming,
		Options: types.MethodOptions{
			IsStreamingResponse: true,
		},
	})

	go func() {
		_ = s.ListenAndServe(ctx)
	}()
	time.Sleep(50 * time.Millisecond)

	client := &http.Client{
		Transport: &qchttp3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://%s/streaming", addr))
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	expected := "Chunk: 0 \n\nChunk: 1 \n\nChunk: 2 \n\nChunk: 3 \n\nChunk: 4 \n\n"
	if strings.ReplaceAll(string(body), "\r", "") != expected {
		t.Fatalf("expected streaming body:\n%q\ngot:\n%q", expected, string(body))
	}
}
