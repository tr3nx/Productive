package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

const pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	var b strings.Builder
	for i := 0; i < l; i++ {
		fmt.Fprintf(&b, "%v", string(pool[randInt(0, len(pool))]))
	}
	return b.String()
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
