package main

import (
	"net/http"

	"github.com/battlej07/goshort/web"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(web.Files))
	mux.HandleFunc("GET /", app.handleHome)
	mux.HandleFunc("POST /shorten", app.handleShorten)
	mux.HandleFunc("GET /{shortenedID}", app.handleRedirect)

	return mux
}
