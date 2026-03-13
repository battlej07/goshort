package main

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/battlej07/goshort/internal/store"
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

const storeTimeout = 3 * time.Second

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, homeTmpl, nil)
}

func (app *application) handleShorten(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var (
		shortened string
		stored    bool
	)
	for range 5 {
		var err error
		shortened, err = generateRandomString()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), storeTimeout)
		err = app.urlStore.Create(ctx, shortened, url)
		cancel()
		if err == nil {
			stored = true
			break
		}

		if errors.Is(err, store.ErrShortURLExists) {
			continue
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !stored {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/result/"+shortened, http.StatusSeeOther)
}

func (app *application) handleResult(w http.ResponseWriter, r *http.Request) {
	shortenedID := r.PathValue("shortenedID")

	ctx, cancel := context.WithTimeout(r.Context(), storeTimeout)
	_, err := app.urlStore.Get(ctx, shortenedID)
	cancel()
	if errors.Is(err, store.ErrShortURLNotFound) {
		w.WriteHeader(http.StatusNotFound)
		renderTemplate(w, notFoundTmpl, nil)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := ResultPageData{getScheme(r) + "://" + r.Host + "/" + shortenedID}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	renderTemplate(w, resultTmpl, data)
}

func (app *application) handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortenedID := r.PathValue("shortenedID")

	ctx, cancel := context.WithTimeout(r.Context(), storeTimeout)
	shortURL, err := app.urlStore.Get(ctx, shortenedID)
	cancel()
	if errors.Is(err, store.ErrShortURLNotFound) {
		w.WriteHeader(http.StatusNotFound)
		renderTemplate(w, notFoundTmpl, nil)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, shortURL.URL, http.StatusFound)
}
