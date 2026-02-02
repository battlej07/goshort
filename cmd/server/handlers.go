package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/battlej07/goshort/web"
)

type ResultPageData struct {
	ShortURL string
}

var (
	homeTmpl     = template.Must(template.ParseFS(web.Files, "views/home.html"))
	resultTmpl   = template.Must(template.ParseFS(web.Files, "views/result.html"))
	notFoundTmpl = template.Must(template.ParseFS(web.Files, "views/not-found.html"))
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if err := homeTmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleShorten(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		url := r.FormValue("url")
		if url == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		shortened, err := generateRandomString()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		db[shortened] = url

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		data := ResultPageData{scheme + "://" + r.Host + "/" + shortened}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := resultTmpl.Execute(w, data); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func handleRedirect(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortenedID := r.PathValue("shortenedID")

		url, ok := db[shortenedID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			if err := notFoundTmpl.Execute(w, nil); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}

func generateRandomString() (string, error) {
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer)[:6], nil
}
