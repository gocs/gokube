package web

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func NewServer(ctx context.Context, addr string) (func(), func() error) {
	srv := &http.Server{
		Addr:        addr,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	return func() {
		webserver(srv)
	}, func() error {
		log.Printf("Shutting down")
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down: %v", err)
		}
		log.Printf("Shutdown complete")
		return nil
	}
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func webserver(srv *http.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greet)
	srv.Handler = mux

	log.Printf("Starting server on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Error starting server: %v", err)
	}
	log.Printf("Server stopped")
}
