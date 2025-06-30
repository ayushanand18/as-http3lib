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

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"

	qchttp3 "github.com/quic-go/quic-go/http3"
)

type DummyResponse struct {
	Key   string `json:"key"`
	Value uint32 `json:"value"`
}

func TestUserRoute_NaiveJSONResponse(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:443"

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	err := server.AddServeMethod(ctx, types.ServeOptions{
		URL:          "/json",
		Method:       constants.HTTP_METHOD_GET,
		ResponseType: constants.RESPONSE_TYPE_JSON_RESPONSE,
		Handler: func(ctx context.Context, r *http.Request) interface{} {
			return DummyResponse{
				Key:   "test",
				Value: 123,
			}
		},
	})
	if err != nil {
		t.Fatalf("AddServeMethod failed: %v", err)
	}

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
