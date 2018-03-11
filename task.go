package main

import (
	"sort"
	"time"
	"fmt"
)

// CREATE TABLE `tasks` (`id` INTEGER PRIMARY KEY AUTO_INCREMENT, `label` VARCHAR(255) NOT NULL, `groupid` INTEGER NOT NULL, `userid` INTEGER NOT NULL, `order` INTEGER NOT NULL, `completed` BIGINT NULL DEFAULT NULL, `created` BIGINT NOT NULL);
var sqlCreateTasksTable = []string{
	"CREATE TABLE `tasks` (",
	"`id` INTEGER PRIMARY KEY AUTO_INCREMENT,",
	"`label` VARCHAR(255) NOT NULL,",
	"`groupid` INTEGER NOT NULL,",
	"`userid` INTEGER NOT NULL,",
	"`order` INTEGER NOT NULL,",
	"`completed` BIGINT NULL DEFAULT NULL,",
	"`created` BIGINT NOT NULL",
	");",
}

type Task struct {
	Id        int    `json:"id"`
	Label     string `json:"label"`
	Groupid   int    `json:"groupid"`
	Userid    int    `json:"userid"`
	Order     int    `json:"order"`
	Completed int64   `json:"completed"`
	Created   int64   `json:"created"`
}

type Tasks []Task

func NewTask(label string, groupid, userid, order int, completed int64) *Task {
	return &Task{
		Label:     label,
		Groupid:   groupid,
		Userid:    userid,
		Order:     order,
		Completed: completed,
		Created:  time.Now().Unix(),
	}
}

func (t *Task) Save() error {
	stmt, err := db.Prepare("INSERT INTO tasks(label, groupid, userid, order, completed, created) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(t.Label, t.Groupid, t.Userid, t.Order, t.Completed, t.Created)
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) UpdateField(field string, change interface{}) error {
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE tasks SET %v=? WHERE id=?", field))
	if err != nil {
		return err
	}
	res, err := stmt.Exec(change, t.Id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Delete() error {
	stmt, err := db.Prepare("DELETE FROM tasks WHERE id=?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(t.Id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (ts Tasks) SortByOrder() {
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Order < ts[j].Order
	})
}

func TasksAll() Tasks {
	var tasks Tasks
	rows, err := db.Query("SELECT id, label, groupid, userid, order, completed, created FROM tasks")
	if err != nil {
		panic(err)
		return tasks
	}
	err = rows.Err()
	if err != nil {
		panic(err)
		return tasks
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Label, &task.Groupid, &task.Userid, &task.Order, &task.Completed, &task.Created)
		if err != nil {
			panic(err)
			return tasks
		}
		tasks = append(tasks, task)
	}
	return tasks
}

func TaskBy(field string, value interface{}) Task {
	var task Task
	err := db.QueryRow(fmt.Sprintf("SELECT id, label, groupid, userid, order, completed, created FROM tasks WHERE %s=?", field), value).Scan(&task.Id, &task.Label, &task.Groupid, &task.Userid, &task.Order, &task.Completed, &task.Created)
	if err != nil {
		panic(err)
		return task
	}
	return task
}

func TasksBy(field string, value interface{}) Tasks {
	var tasks Tasks
	rows, err := db.Query(fmt.Sprintf("SELECT id, label, groupid, userid, order, completed, created FROM tasks WHERE %s=?", field), value)
	if err != nil {
		panic(err)
		return tasks
	}
	err = rows.Err()
	if err != nil {
		panic(err)
		return tasks
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Label, &task.Groupid, &task.Userid, &task.Order, &task.Completed, &task.Created)
		if err != nil {
			panic(err)
			return tasks
		}
		tasks = append(tasks, task)
	}
	return tasks
}
