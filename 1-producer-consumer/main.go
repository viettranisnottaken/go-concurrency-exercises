//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, msq chan<- *Tweet) {
	defer close(msq)

	// because of the way stream.Next is implemented, we cannot parallelize producer
	// If we use goroutines here, we need to use Mutex, which kinda defeats the point
	for {
		tweet, err := stream.Next()
		if errors.Is(err, ErrEOF) {
			return
		}

		msq <- tweet
	}
}

func consumer(msq <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()

	for t := range msq {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func Run() {
	start := time.Now()
	stream := GetMockStream()
	msq := make(chan *Tweet)
	var wg sync.WaitGroup

	// Producer
	go producer(stream, msq)

	// Consumer
	for i := 0; i < 3; i++ { // parallelize consumers
		wg.Add(1)
		go consumer(msq, &wg)
	}

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
