package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data any) {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func generateRandomString() (string, error) {
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer)[:6], nil
}
