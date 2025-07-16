package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

type MyCustomResponseType struct {
	UserId  string
	Message string
}

func main() {
	ctx := context.Background()

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.AddServeMethod(ctx, types.ServeOptions{
		URL: "/users/{user_id}",
		Handler: func(ctx context.Context, r *http.Request) (interface{}, error) {
			// headers := make(map[string]string)
			// headers["X-User-Id"] = r.PathValue("user_id")

			return &MyCustomResponseType{
				UserId:  r.PathValue("user_id"),
				Message: "Hello World from GET.",
			}, nil
		},
		Method: "GET",
	})

	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
