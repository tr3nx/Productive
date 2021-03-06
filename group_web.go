package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func groupHandlers(r *mux.Router) {
	r.HandleFunc("/api/groups", groupsIndex)
	r.HandleFunc("/api/groups/{id:[0-9]+}", groupsSingle)
	r.HandleFunc("/api/groups/create", groupsCreate)
	r.HandleFunc("/api/groups/{id:[0-9]+}/edit", groupsEdit)
	r.HandleFunc("/api/groups/{id:[0-9]+}/delete", groupsDelete)
}

func groupsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	dirtyuserid := r.URL.Query().Get("userid")
	if dirtyuserid != "" {
		userid, err := strconv.Atoi(dirtyuserid)
		if err != nil {
			jsonError(w, err)
			return
		}

		groups, err := GroupsBy("Userid", userid)
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonData(w, groups)
		return
	}
	jsonData(w, GroupsAll())
}

func groupsSingle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	dirtyid := vars["id"]
	id, err := strconv.Atoi(dirtyid)
	if err != nil {
		jsonError(w, err)
		return
	}
	group, err := GroupBy("Id", id)
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, group)
}

func groupsCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	postdata := struct {
		Userid int    `json:"userid"`
		Label  string `json:"label"`
		Order  int    `json:"order"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&postdata)
	if err != nil {
		jsonError(w, err)
		return
	}

	if postdata.Label == "" {
		jsonError(w, errors.New("Label is required"))
		return
	}

	group := NewGroup(postdata.Label, postdata.Userid, postdata.Order)

	err = group.Save()
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, group)
}

func groupsEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var err error
	var newgroup Group
	err = json.NewDecoder(r.Body).Decode(&newgroup)
	if err != nil {
		jsonError(w, err)
		return
	}

	group, err := GroupBy("Id", newgroup.Id)
	if err != nil {
		jsonError(w, err)
		return
	}

	if group.Label != newgroup.Label {
		group.Label = newgroup.Label
		err = group.UpdateField("Label", newgroup.Label)
		if err != nil {
			jsonError(w, err)
			return
		}
	}
	jsonData(w, group)
}

func groupsDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var group Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		jsonError(w, err)
		return
	}

	err = group.Delete()
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, group)
}
