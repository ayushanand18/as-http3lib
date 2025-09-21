package main

import (
	"context"
	"log"
	"os"

	"github.com/ayushanand18/crazyhttp/pkg/errors"
	crazyserver "github.com/ayushanand18/crazyhttp/pkg/server"
)

func main() {
	ctx := context.Background()

	server := crazyserver.NewHttpServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.GET("/audio").Serve(func(ctx context.Context, request interface{}) (interface{}, error) {
		audioBytes, err := os.ReadFile("complete_quest_requirement.mp3")
		if err != nil {
			return nil, errors.InternalServerError.New("Could not read audio file.")
		}

		return audioBytes, nil
	})

	server.GET("/html_file.html").Serve(func(ctx context.Context, request interface{}) (interface{}, error) {
		fileBytes, err := os.ReadFile("html_file.html")
		if err != nil {
			return nil, errors.InternalServerError.New("Could not read html file.")
		}
		return fileBytes, nil
	})

	if err := server.ListenAndServe(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
