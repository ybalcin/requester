package main

import (
	"context"
	"github.com/ybalcin/requester/cmd"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		cancel()
	}()

	if err := cmd.SendRequests(ctx); err != nil {
		panic(err)
	}
}
