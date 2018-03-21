package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func taskHandlers(r *mux.Router) {
	r.HandleFunc("/api/tasks", tasksIndex)
	r.HandleFunc("/api/tasks/{id:[0-9]+}", tasksSingle)
	r.HandleFunc("/api/tasks/create", tasksCreate)
	r.HandleFunc("/api/tasks/{id:[0-9]+}/edit", tasksEdit)
	r.HandleFunc("/api/tasks/{id:[0-9]+}/delete", tasksDelete)
}

func tasksIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	dirtygroupid := r.URL.Query().Get("groupid")
	if dirtygroupid != "" {
		groupid, err := strconv.Atoi(dirtygroupid)
		if err != nil {
			jsonError(w, err)
			return
		}

		tasks, err := TasksBy("groupid", groupid)
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonData(w, tasks)
		return
	}
	jsonData(w, TasksAll())
}

func tasksSingle(w http.ResponseWriter, r *http.Request) {
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
	task, err := TaskBy("Id", id)
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, task)
}

func tasksCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	postdata := struct {
		Label     string `json:"label"`
		Groupid   int    `json:"groupid"`
		Userid    int    `json:"userid"`
		Order     int    `json:"order"`
		Completed bool   `json:"completed"`
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

	var completed bool
	if postdata.Completed {
		completed = true
	}

	task := NewTask(postdata.Label, postdata.Groupid, postdata.Userid, postdata.Order, completed)

	err = task.Save()
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, task)
}

func tasksEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	postdata := struct {
		Label     string `json:"label"`
		Groupid   int    `json:"groupid"`
		Userid    int    `json:"userid"`
		Order     int    `json:"order"`
		Completed bool   `json:"completed"`
	}{}

	var err error
	err = json.NewDecoder(r.Body).Decode(&postdata)
	if err != nil {
		jsonError(w, err)
		return
	}

	vars := mux.Vars(r)
	dirtyid := vars["id"]
	id, err := strconv.Atoi(dirtyid)
	if err != nil {
		jsonError(w, err)
		return
	}

	task, err := TaskBy("Id", id)
	if err != nil {
		jsonError(w, err)
		return
	}

	if task.Userid != postdata.Userid {
		task.Userid = postdata.Userid
		err = task.UpdateField("Userid", postdata.Userid)
		if err != nil {
			jsonError(w, err)
			return
		}
	}

	if task.Groupid != postdata.Groupid {
		task.Groupid = postdata.Groupid
		err = task.UpdateField("Groupid", postdata.Groupid)
		if err != nil {
			jsonError(w, err)
			return
		}
	}

	if task.Label != postdata.Label {
		task.Label = postdata.Label
		err = task.UpdateField("Label", postdata.Label)
		if err != nil {
			jsonError(w, err)
			return
		}
	}

	if postdata.Completed {
		err = task.UpdateField("Completed", task.Completed)
		if err != nil {
			jsonError(w, err)
			return
		}
	}
	jsonData(w, task)
}

func tasksDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		jsonError(w, err)
		return
	}

	err = task.Delete()
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, task)
}
