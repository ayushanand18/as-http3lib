package main_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	crazyserver "github.com/ayushanand18/crazyhttp/pkg/server"
	"github.com/ayushanand18/crazyhttp/pkg/types"

	qchttp3 "github.com/quic-go/quic-go/http3"
)

type DummyResponse struct {
	Key   string `json:"key"`
	Value uint32 `json:"value"`
}

func TestUserRoute_NaiveJSONResponse(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:4431"

	server := crazyserver.NewHttpServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	server.GET("/json").Serve(types.ServeOptions{
		Handler: func(ctx context.Context, req interface{}) (interface{}, error) {
			return DummyResponse{
				Key:   "test",
				Value: 123,
			}, nil
		},
	})

	go func() {
		_ = server.ListenAndServe(ctx)
	}()
	time.Sleep(50 * time.Millisecond)

	client := &http.Client{
		Transport: &qchttp3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://%s/json", addr))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	value := DummyResponse{}

	if err := json.Unmarshal(body, &value); err != nil {
		t.Errorf("Error while parsing JSON, got %q", body)
	}

}
