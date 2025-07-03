package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocs/gokube/internal/web"
	"github.com/gocs/gokube/internal/youtube"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	webserver, shutdown := web.NewServer(ctx, ":8080")
	defer shutdown()

	youtube, err := youtube.NewYoutubeCtrl(ctx)
	if err != nil {
		log.Fatalf("Error creating YouTube controller: %v", err)
	}

	go webserver(map[string]http.HandlerFunc{
		"GET /":                     web.Frontend,
		"GET /api/youtube/playlist": youtube.Playlist,
	})

	<-ctx.Done()
}
