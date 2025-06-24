package main_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
	qchttp3 "github.com/quic-go/quic-go/http3"
)

func TestHTTP3Server_BasicStreamingResponse(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:4434"

	os.Setenv("SERVICE_LISTEN_ADDRESS", addr)

	s := http3.NewServer(ctx)
	if err := s.Initialize(ctx); err != nil {
		t.Fatalf("server initialization failed: %v", err)
	}

	err := s.AddServeMethod(ctx, types.ServeOptions{
		URL:          "/streaming",
		ResponseType: constants.RESPONSE_TYPE_STREAMING_RESPONSE,
		Handler: func(ctx context.Context, r *http.Request) interface{} {
			for i := range 5 {
				time.Sleep(time.Duration(1) * time.Second)

				ctx.Value(constants.STREAMING_RESPONSE_CHANNEL_CONTEXT_KEY).(chan types.StreamChunk) <- types.StreamChunk{
					Id:   uint32(i),
					Data: []byte(fmt.Sprintf("Chunk: %d \n", i)),
				}
			}

			return nil
		},
		Method: "GET",
	})
	if err != nil {
		t.Fatalf("failed to add serve method: %v", err)
	}

	go func() {
		_ = s.ListenAndServe()
	}()
	time.Sleep(50 * time.Millisecond)

	client := &http.Client{
		Transport: &qchttp3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
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

	expected := "Chunk: 0 \nChunk: 1 \nChunk: 2 \nChunk: 3 \nChunk: 4 \n"
	if strings.ReplaceAll(string(body), "\r", "") != expected {
		t.Fatalf("expected streaming body:\n%q\ngot:\n%q", expected, string(body))
	}
}
