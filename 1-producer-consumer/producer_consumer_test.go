package main

import (
	"testing"
	"time"
)

func TestProducerConsumer(t *testing.T) {
	start := time.Now()

	Run()

	if duration := time.Since(start).Seconds(); duration >= 2 {
		t.Errorf("Expected to run under 2 seconds, instead got %f", duration)
	}
}
