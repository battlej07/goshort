package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"

	"github.com/battlej07/goshort/web"
)

type ResultPageData struct {
	ShortURL string
}

var (
	homeTmpl    = template.Must(template.ParseFS(web.Files, "views/home.html"))
	resultTmpl  = template.Must(template.ParseFS(web.Files, "views/result.html"))
	notFundTmpl = template.Must(template.ParseFS(web.Files, "views/not-found.html"))
)

func main() {
	mux := http.NewServeMux()

	db := map[string]string{}

	mux.Handle("GET /static/", http.FileServerFS(web.Files))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_ = homeTmpl.Execute(w, nil)
	})

	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
		}

		url := r.FormValue("url")

		shortend, err := generateRandomString()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		db[shortend] = url

		sheme := "http"
		if r.TLS != nil {
			sheme = "https"
		}

		data := ResultPageData{sheme + "://" + r.Host + "/" + shortend}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		resultTmpl.Execute(w, data)
	})

	mux.HandleFunc("GET /{shortendID}", func(w http.ResponseWriter, r *http.Request) {
		shortendID := r.PathValue("shortendID")

		url, isOk := db[shortendID]
		if !isOk {
			notFundTmpl.Execute(w, nil)
		}

		http.Redirect(w, r, url, http.StatusFound)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Listening on port 8080")
	log.Fatal(srv.ListenAndServe())
}

func generateRandomString() (string, error) {
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer)[:6], nil
}
