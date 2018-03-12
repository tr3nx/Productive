package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"
)

type Task struct {
	Id        int    `json:"id"`
	Label     string `json:"label"`
	Groupid   int    `json:"groupid"`
	Userid    int    `json:"userid"`
	Order     int    `json:"order"`
	Completed int64  `json:"completed"`
	Created   int64  `json:"created"`
}

type Tasks []Task

var taskfields = []string{"id", "label", "groupid", "userid", "order", "completed", "created"}

func init() {
	dbRegisterMigration("TasksCreateTable", TasksCreateTable)
	log.Println("[#] Tasks module loading...")
}

func NewTask(label string, groupid, userid, order int, completed int64) *Task {
	return &Task{
		Label:     label,
		Groupid:   groupid,
		Userid:    userid,
		Order:     order,
		Completed: completed,
		Created:   time.Now().Unix(),
	}
}

func (t *Task) Save() error {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO `tasks`(%v) values(?, ?, ?, ?, ?, ?)", joinFields(taskfields[1:])))
	if err != nil {
		return err
	}
	res, err := stmt.Exec(t.Label, t.Groupid, t.Userid, t.Order, t.Completed, t.Created)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.Id = int(id)
	return nil
}

func (t *Task) UpdateField(field string, change interface{}) error {
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE `tasks` SET `%v`=? WHERE `id`=?", field))
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
	stmt, err := db.Prepare("DELETE FROM `tasks` WHERE `id`=?")
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

func TasksCreateTable() error {
	stmt, err := db.Prepare("CREATE TABLE `tasks` (`id` INTEGER PRIMARY KEY AUTO_INCREMENT, `label` VARCHAR(255) NOT NULL, `groupid` INTEGER NOT NULL, `userid` INTEGER NOT NULL, `order` INTEGER NOT NULL, `completed` BIGINT NULL DEFAULT NULL, `created` BIGINT NOT NULL)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec()
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func TasksDropTable() error {
	stmt, err := db.Prepare("DROP TABLE `tasks`")
	if err != nil {
		return err
	}
	res, err := stmt.Exec()
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
	rows, err := db.Query(fmt.Sprintf("SELECT %v FROM `tasks`", joinFields(taskfields)))
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

func TaskBy(field string, value interface{}) (Task, error) {
	var task Task
	err := db.QueryRow(fmt.Sprintf("SELECT %v FROM `tasks` WHERE `%v`=?", joinFields(taskfields), field), value).Scan(&task.Id, &task.Label, &task.Groupid, &task.Userid, &task.Order, &task.Completed, &task.Created)
	if err != nil {
		return task, err
	}
	if reflect.DeepEqual(task, Task{}) {
		return task, errors.New("No task found")
	}
	return task, nil
}

func TasksBy(field string, value interface{}) (Tasks, error) {
	var tasks Tasks
	rows, err := db.Query(fmt.Sprintf("SELECT %v FROM `tasks` WHERE `%v`=?", joinFields(taskfields), field), value)
	if err != nil {
		return tasks, err
	}
	err = rows.Err()
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Label, &task.Groupid, &task.Userid, &task.Order, &task.Completed, &task.Created)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	if len(tasks) <= 0 {
		return tasks, errors.New("No tasks found")
	}
	return tasks, nil
}
