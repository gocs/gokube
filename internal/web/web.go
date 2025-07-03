package web

import (
	"context"
	"log"
	"net"
	"net/http"
)

func NewServer(ctx context.Context, addr string) (func(handlers map[string]http.HandlerFunc), func() error) {
	srv := &http.Server{
		Addr:        addr,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	return func(handlers map[string]http.HandlerFunc) {
			webserver(srv, handlers)
		}, func() error {
			log.Printf("Shutting down")
			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down: %v", err)
			}
			log.Printf("Shutdown complete")
			return nil
		}
}

func webserver(srv *http.Server, handlers map[string]http.HandlerFunc) {
	mux := http.NewServeMux()
	for path, handler := range handlers {
		mux.HandleFunc(path, handler)
	}
	srv.Handler = mux

	log.Printf("Starting server on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Error starting server: %v", err)
	}
	log.Printf("Server stopped")
}
