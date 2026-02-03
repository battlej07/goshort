package main

import (
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

	shortened, err := generateRandomString()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	app.db[shortened] = url

	http.Redirect(w, r, "/result/"+shortened, http.StatusSeeOther)
}

func (app *application) handleResult(w http.ResponseWriter, r *http.Request) {
	shortenedID := r.PathValue("shortenedID")

	_, ok := app.db[shortenedID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		renderTemplate(w, notFoundTmpl, nil)
		return
	}

	data := ResultPageData{getScheme(r) + "://" + r.Host + "/" + shortenedID}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	renderTemplate(w, resultTmpl, data)
}

func (app *application) handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortenedID := r.PathValue("shortenedID")

	url, ok := app.db[shortenedID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		renderTemplate(w, notFoundTmpl, nil)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
