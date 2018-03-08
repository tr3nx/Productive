package main

import (
	"github.com/asdine/storm/q"
	"sort"
)

type Task struct {
	Id        int    `json:"id" storm:"id,increment"`
	Groupid   int    `json:"groupid"`
	Userid    int    `json:"userid"`
	Order     int    `json:"order"`
	Label     string `json:"label"`
	Completed bool   `json:"completed,bool"`
}

type Tasks []Task

func NewTask(groupid, userid, order int, label string, completed bool) *Task {
	return &Task{
		Groupid:   groupid,
		Userid:    userid,
		Order:     order,
		Label:     label,
		Completed: completed,
	}
}

func (t *Task) Save() error {
	return db.Save(t)
}

func (t *Task) UpdateField(field string, change interface{}) error {
	return db.UpdateField(t, field, change)
}

func (t *Task) Delete() error {
	return db.DeleteStruct(t)
}

func (ts Tasks) SortByOrder() {
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Order < ts[j].Order
	})
}

func TasksAll() Tasks {
	var tasks Tasks
	err := db.All(&tasks)
	if err != nil {
		panic(err)
		return tasks
	}
	return tasks
}

func TaskBy(field string, value interface{}) Task {
	var task Task
	err := db.One(field, value, &task)
	if err != nil {
		panic(err)
		return task
	}
	return task
}

func TasksBy(field string, value interface{}) Tasks {
	var tasks Tasks
	err := db.Find(field, value, &tasks)
	if err != nil {
		panic(err)
		return tasks
	}
	return tasks
}

func TasksByUserIdAndCompleted(userid int) Tasks {
	var tasks Tasks
	err := db.Select(q.Eq("Userid", userid), q.Eq("Completed", true)).Find(&tasks)
	if err != nil {
		panic(err)
		return tasks
	}
	return tasks
}

func TasksByUserIdAndNotCompleted(userid int) Tasks {
	var tasks Tasks
	err := db.Select(q.Eq("Userid", userid), q.Eq("Completed", false)).Find(&tasks)
	if err != nil {
		panic(err)
		return tasks
	}
	return tasks
}
