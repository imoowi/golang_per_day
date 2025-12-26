package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	const (
		url         = "http://localhost:8080/order/1"
		concurrency = 20  // 并发数
		requests    = 100 // 总请求数
	)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var success, fail int

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requests/concurrency; j++ {
				resp, err := client.Get(url)
				if err != nil {
					mu.Lock()
					fail++
					mu.Unlock()
					continue
				}
				_, err = io.ReadAll(resp.Body)
				resp.Body.Close()
				mu.Lock()
				if err != nil || resp.StatusCode != 200 {
					fail++
				} else {
					success++
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("=== Stress Test Result ===")
	fmt.Printf("Total Requests: %d\n", requests)
	fmt.Printf("Success: %d\n", success)
	fmt.Printf("Fail (timeout / circuit open): %d\n", fail)
	fmt.Printf("Total Duration: %v\n", duration)
}
