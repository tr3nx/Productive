package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func authHandlers(r *mux.Router) {
	r.HandleFunc("/api/auth/login", authLogin)
	r.HandleFunc("/api/auth/signup", authSignup)
}

func authLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		jsonError(w, err)
		return
	}

	user, err := creds.Authenticate()
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, user)
}

func authSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		jsonError(w, err)
		return
	}

	user, err := creds.Signup()
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, user)
}
