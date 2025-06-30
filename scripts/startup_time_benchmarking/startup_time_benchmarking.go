package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	const TARGET_URL = "https://localhost:443/test"
	const MAX_TRIES = 10000
	// choosing 8ms as context deadline, since p99 < 8ms
	// we want to only wait until request lifetime
	// to get an accurate estimate of startup time
	const TRY_DURATION = time.Duration(8) * time.Millisecond

	lastUnsuccessfullRequest := time.Now()
	firstSuccessfulRequest := time.Now()

	tr := &http3.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer tr.Close()

	client := &http.Client{
		Transport: tr,
	}

	for try := 0; try < MAX_TRIES; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), TRY_DURATION)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, "GET", TARGET_URL, nil)
		_, err := client.Do(req)
		if err != nil {
			fmt.Println("Unsuccessfull")
			lastUnsuccessfullRequest = time.Now()
		} else {
			fmt.Println("Successfull")
			firstSuccessfulRequest = time.Now()
			break
		}
	}

	fmt.Printf("Last Unsuccessfull Request: %+v \n", lastUnsuccessfullRequest)
	fmt.Printf("First Successfull Request: %+v \n", firstSuccessfulRequest)
	fmt.Printf("Time Diff: %+v \n", firstSuccessfulRequest.Sub(lastUnsuccessfullRequest))
}
