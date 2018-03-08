package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleHttp() {
	r := mux.NewRouter()

	taskHandlers(r)
	groupHandlers(r)
	authHandlers(r)

	r.HandleFunc("/", appHandler)

	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/")))
	r.PathPrefix("/assets/").Handler(fs)

	log.Println("[@] Http listening...")

	log.Fatal(http.ListenAndServe(":5100", r))
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	htmlResponse(w, renderTemplate("layouts/app.html", nil))
}
