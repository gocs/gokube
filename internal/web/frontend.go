package web

import (
	"embed"
	"net/http"
)

//go:embed srv
var srv embed.FS

func Frontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, srv, "srv/index.html")
}
