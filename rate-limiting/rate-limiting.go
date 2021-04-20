package main

import (
	"fmt"
	"time"
)

func main() {
	requests := buildRequestChan(5)

	limiter := time.Tick(500 * time.Millisecond)

	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	burstyRequests := buildRequestChan(5)

	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(500 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("bursty request", req, time.Now())
	}
}

func buildRequestChan(size int) chan int {
	requestChan := make(chan int, 5)

	defer close(requestChan)

	for i := 1; i <= 5; i++ {
		requestChan <- i
	}

	return requestChan
}
