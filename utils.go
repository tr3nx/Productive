package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

func HashPassword(pass string) string {
	return pass
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(data)
}

func jsonError(w http.ResponseWriter, err error) {
	jsonResponse(w, struct {
		Success bool   `json:"success,bool"`
		Message string `json:"message"`
	}{
		Success: false,
		Message: err.Error(),
	})
}

func jsonData(w http.ResponseWriter, data interface{}) {
	jsonResponse(w, struct {
		Success bool        `json:"success,bool"`
		Data    interface{} `json:"data"`
	}{
		Success: true,
		Data:    data,
	})
}

func htmlResponse(w http.ResponseWriter, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, html)
}

func renderTemplate(file string, data interface{}) string {
	name := strings.Split(file, "/")
	b := new(bytes.Buffer)
	tpl := template.Must(template.New(name[len(name)-1]).ParseFiles(file))
	tpl.Execute(b, data)
	return b.String()
}

func deleteIfExists(path string) error {
	if fileExists(path) {
		return os.Remove(path)
	}
	return nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
