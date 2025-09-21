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

	"github.com/ayushanand18/crazyhttp/pkg/errors"
	crazyserver "github.com/ayushanand18/crazyhttp/pkg/server"
	"github.com/ayushanand18/crazyhttp/pkg/types"
	qchttp3 "github.com/quic-go/quic-go/http3"
)

func setupTestServer(ctx context.Context) (*crazyserver.HttpServer, string) {
	addr := "localhost:4431"
	server := crazyserver.NewHttpServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		panic(fmt.Sprintf("Server initialization failed: %v", err))
	}

	server.GET("/audio").Serve(types.ServeOptions{
		Handler: func(ctx context.Context, request interface{}) (interface{}, error) {
			audioBytes, err := os.ReadFile("complete_quest_requirement.mp3")
			if err != nil {
				return nil, errors.InternalServerError.New("Could not read audio file.")
			}
			return audioBytes, nil
		},
	})

	server.GET("/html_file.html").Serve(types.ServeOptions{
		Handler: func(ctx context.Context, request interface{}) (interface{}, error) {
			fileBytes, err := os.ReadFile("html_file.html")
			if err != nil {
				return nil, errors.InternalServerError.New("Could not read html file.")
			}
			return fileBytes, nil
		},
	})

	go func() {
		_ = server.ListenAndServe(ctx)
	}()
	time.Sleep(100 * time.Millisecond) // Give server time to start
	return &server, addr
}

func createClient() *http.Client {
	return &http.Client{
		Transport: &qchttp3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func TestAudioEndpoint(t *testing.T) {
	ctx := context.Background()
	_, addr := setupTestServer(ctx)

	client := createClient()
	resp, err := client.Get(fmt.Sprintf("https://%s/audio", addr))
	if err != nil {
		t.Fatalf("GET /audio failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	if len(body) == 0 {
		t.Errorf("Expected non-empty audio file body")
	}
}

func TestHTMLEndpoint(t *testing.T) {
	ctx := context.Background()
	_, addr := setupTestServer(ctx)

	client := createClient()
	resp, err := client.Get(fmt.Sprintf("https://%s/html_file.html", addr))
	if err != nil {
		t.Fatalf("GET /html_file.html failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	if len(body) == 0 {
		t.Errorf("Expected non-empty HTML file body")
	}
}
