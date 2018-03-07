package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func taskHandlers(r *mux.Router) {
	r.HandleFunc("/api/tasks", tasksIndex)
	r.HandleFunc("/api/tasks/{id:[0-9]+}", tasksSingle)
	r.HandleFunc("/api/tasks/create", tasksCreate)
	r.HandleFunc("/api/tasks/{id:[0-9]+}/edit", tasksEdit)
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

		var tasks Tasks
		err = db.Find("Groupid", groupid, &tasks)
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
	jsonData(w, TaskBy("Id", id))
}

func tasksCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	postdata := struct {
		Userid    int    `json:"userid"`
		Groupid   int    `json:"groupid"`
		Order     int    `json:"order"`
		Label     string `json:"label"`
		Completed bool   `json:"completed,bool"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&postdata)
	if err != nil {
		jsonError(w, err)
		return
	}

	task := NewTask(postdata.Groupid, postdata.Userid, postdata.Order, postdata.Label, postdata.Completed)

	err = task.Save()
	if err != nil {
		jsonError(w, err)
		return
	}
	jsonData(w, *task)
}

func tasksEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - Method is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var err error
	var newtask Task
	err = json.NewDecoder(r.Body).Decode(&newtask)
	if err != nil {
		jsonError(w, err)
		return
	}

	task := TaskBy("Id", newtask.Id)

	if task.Userid != newtask.Userid {
		task.Userid = newtask.Userid
		err = task.UpdateField("Userid", newtask.Userid)
		if err != nil {
			jsonError(w, err)
			return
		}
	}

	if task.Groupid != newtask.Groupid {
		task.Groupid = newtask.Groupid
		err = task.UpdateField("Groupid", newtask.Groupid)
		if err != nil {
			jsonError(w, err)
			return
		}
	}

	if task.Label != newtask.Label {
		task.Label = newtask.Label
		err = task.UpdateField("Label", newtask.Label)
		if err != nil {
			jsonError(w, err)
			return
		}
	}

	if task.Completed != newtask.Completed {
		task.Completed = newtask.Completed
		err = task.UpdateField("Completed", newtask.Completed)
		if err != nil {
			jsonError(w, err)
			return
		}
	}
	jsonData(w, *task)
}
