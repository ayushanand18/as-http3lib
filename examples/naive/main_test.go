package main_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
	qchttp3 "github.com/quic-go/quic-go/http3"
)

func TestUserRoute_NaiveGETRequest(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:4433"
	os.Setenv("SERVICE_LISTEN_ADDRESS", addr)

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	err := server.AddServeMethod(ctx, types.ServeOptions{
		URL:          "/test",
		Method:       constants.HTTP_METHOD_GET,
		ResponseType: constants.RESPONSE_TYPE_BASE_RESPONSE,
		Handler: func(ctx context.Context, r *http.Request) interface{} {
			headers := map[string]string{
				"X-User-Id": r.PathValue("user_id"),
			}
			return &types.HttpResponse{
				StatusCode: 200,
				Headers:    headers,
				Body:       []byte("Hello World from GET."),
			}
		},
	})
	if err != nil {
		t.Fatalf("AddServeMethod failed: %v", err)
	}

	err = server.AddServeMethod(ctx, types.ServeOptions{
		URL: "/test",
		Handler: func(ctx context.Context, r *http.Request) interface{} {
			return &types.HttpResponse{
				StatusCode: 200,
				Body:       []byte("Hello World from POST."),
			}
		},
		Method: "POST",
	})
	if err != nil {
		t.Fatalf("AddServeMethod failed: %v", err)
	}

	go func() {
		_ = server.ListenAndServe()
	}()
	time.Sleep(50 * time.Millisecond)

	client := &http.Client{
		Transport: &qchttp3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://%s/test", addr))
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
	if string(body) != "Hello World from GET." {
		t.Errorf("Expected body 'Hello World from GET.', got %q", body)
	}

}

func TestUserRoute_NaivePOSTRequest(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:4433"
	os.Setenv("SERVICE_LISTEN_ADDRESS", addr)

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	err := server.AddServeMethod(ctx, types.ServeOptions{
		URL:          "/test",
		Method:       constants.HTTP_METHOD_GET,
		ResponseType: constants.RESPONSE_TYPE_BASE_RESPONSE,
		Handler: func(ctx context.Context, r *http.Request) interface{} {
			return &types.HttpResponse{
				StatusCode: 200,
				Body:       []byte("Hello World from POST."),
			}
		},
	})
	if err != nil {
		t.Fatalf("AddServeMethod failed: %v", err)
	}

	err = server.AddServeMethod(ctx, types.ServeOptions{
		URL: "/test",
		Handler: func(ctx context.Context, r *http.Request) interface{} {
			return &types.HttpResponse{
				StatusCode: 200,
				Body:       []byte("Hello World from POST."),
			}
		},
		Method: "POST",
	})
	if err != nil {
		t.Fatalf("AddServeMethod failed: %v", err)
	}

	go func() {
		_ = server.ListenAndServe()
	}()
	time.Sleep(50 * time.Millisecond)

	client := &http.Client{
		Transport: &qchttp3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	contentType := "text/plain; utf-8"
	resp, err := client.Post(fmt.Sprintf("https://%s/test", addr), contentType, nil)
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
	if string(body) != "Hello World from POST." {
		t.Errorf("Expected body 'Hello World from POST.', got %q", body)
	}

}
