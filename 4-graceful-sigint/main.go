//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func gracefulShutdown(proc *MockProcess, done chan<- struct{}) {
	sigint, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	<-sigint.Done()

	log.Println("Shutting down gracefully")

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Println("Waiting 5 seconds for things to settle down")
	<-timeout.Done()

	proc.Stop()

	log.Println("Shut down")
	done <- struct{}{}
}

func main() {
	// Create a process
	proc := MockProcess{}
	done := make(chan struct{}, 1)
	go gracefulShutdown(&proc, done)

	// Run the process (blocking)
	proc.Run()
	<-done
}
