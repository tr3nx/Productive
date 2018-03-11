package main

import (
	"sort"
	"time"
	"fmt"
)

// CREATE TABLE `groups` (`id` INTEGER PRIMARY KEY AUTO_INCREMENT, `label` VARCHAR(128) NOT NULL, `userid` INTEGER NOT NULL, `order` INTEGER NOT NULL, `created` BIGINT NOT NULL);
var sqlCreateGroupsTable = []string{
	"CREATE TABLE `groups` (",
	"`id` INTEGER PRIMARY KEY AUTO_INCREMENT,",
	"`label` VARCHAR(128) NOT NULL,",
	"`userid` INTEGER NOT NULL,",
	"`order` INTEGER NOT NULL,",
	"`created` BIGINT NOT NULL",
	");",
}

type Query struct {
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
	Order  string `json:"order"`
	Sort   string `json:"sort"`
}

type Group struct {
	Id     int    `json:"id"`
	Label  string `json:"label"`
	Userid int    `json:"userid"`
	Order  int    `json:"order"`
	Created int64 `json:"created"`
}

type Groups []Group

func NewGroup(label string, userid, order int) *Group {
	return &Group{
		Label:  label,
		Userid: userid,
		Order:  order,
		Created: time.Now().Unix(),
	}
}

func (g *Group) Save() error {
	stmt, err := db.Prepare("INSERT INTO groups(label, userid, order, created) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(g.Label, g.Userid, g.Order, g.Created)
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (g *Group) UpdateField(field string, change interface{}) error {
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE groups SET %v=? WHERE id=?", field))
	if err != nil {
		return err
	}
	res, err := stmt.Exec(change, g.Id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (g *Group) Delete() error {
	stmt, err := db.Prepare("DELETE FROM groups WHERE id=?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(g.Id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (gs Groups) SortByOrder() {
	sort.Slice(gs, func(i, j int) bool {
		return gs[i].Order < gs[j].Order
	})
}

func GroupsAll() Groups {
	var groups Groups
	rows, err := db.Query("SELECT id, label, userid, order, created FROM groups")
	if err != nil {
		panic(err)
		return groups
	}
	err = rows.Err()
	if err != nil {
		panic(err)
		return groups
	}
	defer rows.Close()
	for rows.Next() {
		var group Group
		err = rows.Scan(&group.Id, &group.Label, &group.Userid, &group.Order, &group.Created)
		if err != nil {
			panic(err)
			return groups
		}
		groups = append(groups, group)
	}
	return groups
}

func GroupBy(field string, value interface{}) Group {
	var group Group
	err := db.QueryRow(fmt.Sprintf("SELECT id, label, userid, order, created FROM groups WHERE %s=?", field), value).Scan(&group.Id, &group.Label, &group.Userid, &group.Order, &group.Created)
	if err != nil {
		panic(err)
		return group
	}
	return group
}

func GroupsBy(field string, value interface{}) Groups {
	var groups Groups
	rows, err := db.Query(fmt.Sprintf("SELECT id, label, userid, order, created FROM groups WHERE %s=?", field), value)
	if err != nil {
		panic(err)
		return groups
	}
	err = rows.Err()
	if err != nil {
		panic(err)
		return groups
	}
	defer rows.Close()
	for rows.Next() {
		var group Group
		err = rows.Scan(&group.Id, &group.Label, &group.Userid, &group.Order, &group.Created)
		if err != nil {
			panic(err)
			return groups
		}
		groups = append(groups, group)
	}
	return groups
}
