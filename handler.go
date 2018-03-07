package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleHttp() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)

	taskHandlers(r)
	groupHandlers(r)
	authHandlers(r)

	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/")))
	r.PathPrefix("/assets/").Handler(fs)

	log.Fatal(http.ListenAndServe(":5100", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	htmlResponse(w, renderTemplate("pages/home.html", nil))
}
