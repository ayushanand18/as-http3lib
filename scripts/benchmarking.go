package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go/http3"
)

type RequestStats struct {
	Duration time.Duration
	Success  bool
	Err      error
}

var (
	totalRequests      int64
	successfulRequests int64
	failedRequests     int64
	minLatency         atomic.Value
	maxLatency         atomic.Value
	totalLatencyNs     int64

	allLatencies   []time.Duration
	latenciesMutex sync.Mutex
)

func init() {
	minLatency.Store(time.Duration(1<<63 - 1))
	maxLatency.Store(time.Duration(0))
}

func worker(
	ctx context.Context,
	url string,
	client *http.Client,
	resultsChan chan<- RequestStats,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			startTime := time.Now()
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				select {
				case resultsChan <- RequestStats{Err: fmt.Errorf("failed to create request: %w", err)}:
				case <-ctx.Done():
				}
				continue
			}

			resp, err := client.Do(req)
			if err != nil {
				select {
				case resultsChan <- RequestStats{Err: fmt.Errorf("request failed: %w", err)}:
				case <-ctx.Done():
				}
				continue
			}
			defer resp.Body.Close()

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				select {
				case resultsChan <- RequestStats{Err: fmt.Errorf("failed to read response body: %w", err)}:
				case <-ctx.Done():
				}
				continue
			}

			duration := time.Since(startTime)
			success := resp.StatusCode == http.StatusOK

			select {
			case resultsChan <- RequestStats{
				Duration: duration,
				Success:  success,
				Err:      nil,
			}:
			case <-ctx.Done():
			}

			time.Sleep(10 * time.Millisecond)
		}
	}
}

func calculatePercentile(latencies []time.Duration, percentile float64) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	if percentile <= 0 {
		return latencies[0]
	}
	if percentile >= 100 {
		return latencies[len(latencies)-1]
	}

	index := (percentile / 100.0) * float64(len(latencies)-1)
	if index == float64(int(index)) {
		return latencies[int(index)]
	}
	lowerIndex := int(index)
	upperIndex := lowerIndex + 1
	weight := index - float64(lowerIndex)
	return time.Duration(float64(latencies[lowerIndex])*(1.0-weight) + float64(latencies[upperIndex])*weight)
}

func main() {
	const (
		targetURL    = "https://localhost:4433/"
		virtualUsers = 50
		testDuration = 30 * time.Second
	)

	fmt.Printf("Starting HTTP/3 Load Test on %s\n", targetURL)
	fmt.Printf("Virtual Users: %d, Duration: %s\n", virtualUsers, testDuration)
	fmt.Println("--------------------------------------------------")

	tr := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer tr.Close()

	client := &http.Client{
		Transport: tr,
	}

	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()

	resultsChan := make(chan RequestStats, virtualUsers*100)

	var wg sync.WaitGroup

	for i := 0; i < virtualUsers; i++ {
		wg.Add(1)
		go worker(ctx, targetURL, client, resultsChan, &wg)
	}

	var processingWg sync.WaitGroup
	processingWg.Add(1)
	go func() {
		defer processingWg.Done()
		for stats := range resultsChan {
			atomic.AddInt64(&totalRequests, 1)
			if stats.Success {
				atomic.AddInt64(&successfulRequests, 1)
			} else {
				atomic.AddInt64(&failedRequests, 1)
				if stats.Err != nil {
					log.Printf("Request failed: %v", stats.Err)
				}
			}

			durationNs := stats.Duration.Nanoseconds()
			atomic.AddInt64(&totalLatencyNs, durationNs)

			for {
				oldMin := minLatency.Load().(time.Duration)
				if stats.Duration < oldMin {
					if minLatency.CompareAndSwap(oldMin, stats.Duration) {
						break
					}
				} else {
					break
				}
			}

			for {
				oldMax := maxLatency.Load().(time.Duration)
				if stats.Duration > oldMax {
					if maxLatency.CompareAndSwap(oldMax, stats.Duration) {
						break
					}
				} else {
					break
				}
			}

			latenciesMutex.Lock()
			allLatencies = append(allLatencies, stats.Duration)
			latenciesMutex.Unlock()
		}
	}()

	<-ctx.Done()
	log.Println("Test duration elapsed. Signaling workers to stop...")

	wg.Wait()
	close(resultsChan)
	processingWg.Wait()

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("Load Test Summary:")
	fmt.Printf("Total Duration: %.2f seconds\n", testDuration.Seconds())
	fmt.Printf("Total Requests: %d\n", atomic.LoadInt64(&totalRequests))
	fmt.Printf("Successful Requests: %d\n", atomic.LoadInt64(&successfulRequests))
	fmt.Printf("Failed Requests: %d\n", atomic.LoadInt64(&failedRequests))

	if atomic.LoadInt64(&totalRequests) > 0 {
		rps := float64(atomic.LoadInt64(&totalRequests)) / testDuration.Seconds()
		fmt.Printf("Requests per Second (RPS): %.2f\n", rps)

		avgLatencyNs := float64(atomic.LoadInt64(&totalLatencyNs)) / float64(atomic.LoadInt64(&totalRequests))
		fmt.Printf("Average Latency: %s\n", time.Duration(avgLatencyNs).Round(time.Microsecond))
		fmt.Printf("Min Latency: %s\n", minLatency.Load().(time.Duration).Round(time.Microsecond))
		fmt.Printf("Max Latency: %s\n", maxLatency.Load().(time.Duration).Round(time.Microsecond))

		latenciesMutex.Lock()
		sort.Slice(allLatencies, func(i, j int) bool {
			return allLatencies[i] < allLatencies[j]
		})
		latenciesMutex.Unlock()

		fmt.Println("\nLatency Percentiles:")
		p50 := calculatePercentile(allLatencies, 50)
		p90 := calculatePercentile(allLatencies, 90)
		p95 := calculatePercentile(allLatencies, 95)
		p99 := calculatePercentile(allLatencies, 99)

		fmt.Printf("p50 (Median): %s\n", p50.Round(time.Microsecond))
		fmt.Printf("p90: %s\n", p90.Round(time.Microsecond))
		fmt.Printf("p95: %s\n", p95.Round(time.Microsecond))
		fmt.Printf("p99: %s\n", p99.Round(time.Microsecond))

	} else {
		fmt.Println("No requests were sent during the test.")
	}
	fmt.Println("--------------------------------------------------")
}
