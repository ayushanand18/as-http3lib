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

type MyCustomResponseType struct {
	UserId  string
	Message string
}

func Encoder(ctx context.Context, response interface{}) (headers map[string][]string, body []byte, err error) {
	resp := response.(MyCustomResponseType)
	headers = make(map[string][]string)
	headers["X-User-Id"] = []string{resp.UserId}

	bodyBytes, err := json.Marshal(resp)
	if err != nil {
		return headers, body, err
	}

	return headers, bodyBytes, nil
}

func UserIdHandler(ctx context.Context, request interface{}) (response interface{}, err error) {
	pathValues := ctx.Value(constants.HTTP_REQUEST_PATH_VALUES).(map[string]string)

	return MyCustomResponseType{
		UserId:  pathValues["user_id"],
		Message: "Hello World from GET.",
	}, nil
}

func TestUserRoute_WithUserIdHeader(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:4431"

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	server.GET("/users/{user_id}").Serve(types.ServeOptions{
		Handler: UserIdHandler,
		Encoder: Encoder,
	})

	go func() {
		_ = server.ListenAndServe(ctx)
	}()
	time.Sleep(500 * time.Millisecond)

	client := &http.Client{
		Transport: &qchttp3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://%s/users/123", addr))
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
	v := MyCustomResponseType{}
	if err := json.Unmarshal(body, &v); err != nil {
		t.Errorf("Could not unmarshal json")
	}
	if v.Message != "Hello World from GET." || v.UserId != "123" {
		t.Errorf("Expected Message 'Hello World from GET.' and UserId '123', got %q", body)
	}

	userID := resp.Header.Get("X-User-Id")
	if userID != "123" {
		t.Errorf("Expected X-User-Id=123, got %q", userID)
	}
}
