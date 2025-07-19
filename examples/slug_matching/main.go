package main

import (
	"context"
	"log"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/http3"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

type MyCustomResponseType struct {
	UserId  string
	Message string
}

func UserIdHandler(ctx context.Context, request interface{}) (response interface{}, err error) {
	pathValues := ctx.Value(constants.HTTP_REQUEST_PATH_VALUES).(map[string]string)

	return &MyCustomResponseType{
		UserId:  pathValues["user_id"],
		Message: "Hello World from GET.",
	}, nil
}

func main() {
	ctx := context.Background()

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.GET("/users/{user_id}").Serve(types.ServeOptions{
		Handler: UserIdHandler,
	})

	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
