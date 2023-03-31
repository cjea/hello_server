package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Parse environment variables
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}

	numRequestsStr := os.Getenv("NUM_REQUESTS")
	if numRequestsStr == "" {
		numRequestsStr = "10"
	}

	numWorkersStr := os.Getenv("NUM_WORKERS")
	if numWorkersStr == "" {
		numWorkersStr = "10"
	}

	numRequests, err := strconv.Atoi(numRequestsStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid NUM_REQUESTS: %s", numRequestsStr))
	}

	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid NUM_WORKERS: %s", numWorkersStr))
	}

	// Start worker nodes
	var wg sync.WaitGroup
	requests := make(chan int, numRequests)
	statusCodes := make(chan int, numRequests)

	start := time.Now()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(requests, statusCodes, name, &wg)
	}

	// Send requests
	for i := 0; i < numRequests; i++ {
		requests <- i
	}
	close(requests)

	// Wait for worker nodes to finish
	wg.Wait()
	close(statusCodes)

	elapsed := time.Since(start)
	// Create report
	statusCounts := make(map[int]int)
	for code := range statusCodes {
		statusCounts[code]++
	}

	fmt.Printf("Report:\n")
	fmt.Printf("Total requests: %d\n", numRequests)
	for code, count := range statusCounts {
		fmt.Printf("%d: %d\n", code, count)
	}

	fmt.Printf("HTTP requests took %v seconds\n", elapsed.Seconds())
}

func worker(requests <-chan int, statusCodes chan<- int, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("http://localhost:8080/hello?name=%s", name)

	for range requests {
		resp, err := http.Get(url)
		if err != nil {
			statusCodes <- 0
		} else {
			statusCodes <- resp.StatusCode
			resp.Body.Close()
		}
	}
}
