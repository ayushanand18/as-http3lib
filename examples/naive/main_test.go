package main_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	crazyserver "github.com/ayushanand18/crazyhttp/pkg/server"
	"github.com/ayushanand18/crazyhttp/pkg/types"
)

func TestUserRoute_NaiveGETRequest(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:4431"

	server := crazyserver.NewHttpServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	server.GET("/test").Serve(types.ServeOptions{
		Handler: func(ctx context.Context, request interface{}) (interface{}, error) {
			return "Hello World from GET.", nil
		},
	})

	go func() {
		_ = server.ListenAndServe(ctx)
	}()
	time.Sleep(50 * time.Millisecond)

	client := &http.Client{}

	resp, err := client.Get(fmt.Sprintf("http://%s/test", addr))
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
		t.Errorf("Expected body \"Hello World from GET.\" got %q", body)
	}

}

func TestUserRoute_NaivePOSTRequest(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:4431"

	server := crazyserver.NewHttpServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	server.POST("/test").Serve(types.ServeOptions{
		Handler: func(ctx context.Context, request interface{}) (interface{}, error) {
			return "Hello World from POST.", nil
		},
	})

	go func() {
		_ = server.ListenAndServe(ctx)
	}()
	time.Sleep(50 * time.Millisecond)

	client := &http.Client{}

	contentType := "text/plain; utf-8"
	resp, err := client.Post(fmt.Sprintf("http://%s/test", addr), contentType, nil)
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
