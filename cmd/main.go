package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocs/gokube/internal/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	webserver, shutdown := web.NewServer(ctx, ":8080")
	defer shutdown()

	go webserver()

	<-ctx.Done()
}
