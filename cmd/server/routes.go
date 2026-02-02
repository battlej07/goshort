package main

import (
	"net/http"

	"github.com/battlej07/goshort/web"
)

func registerRoutes(mux *http.ServeMux, db map[string]string) {
	mux.Handle("GET /static/", http.FileServerFS(web.Files))
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("POST /shorten", handleShorten(db))
	mux.HandleFunc("GET /{shortenedID}", handleRedirect(db))
}
