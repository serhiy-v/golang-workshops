package main

import (
	"fmt"
	"time"
)

func producer(stream Stream, ch chan *Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(ch)
			return
		}
		ch <- tweet
	}
}

func consumer(tweets chan *Tweet) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	ch := make(chan *Tweet)
	// Hint: this can be resolved via channels
	// Producer
	go producer(stream, ch)
	// Consumer
	consumer(ch)

	fmt.Printf("Process took %s\n", time.Since(start))
}
