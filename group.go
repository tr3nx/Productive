package main

import (
	"sort"
)

type Query struct {
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
	Order  string `json:"order"`
	Sort   string `json:"sort"`
}

type Group struct {
	Id     int    `json:"id" storm:"id,increment"`
	Userid int    `json:"userid"`
	Order  int    `json:"order"`
	Label  string `json:"label"`
}

type Groups []Group

func NewGroup(userid, order int, label string) *Group {
	return &Group{
		Userid: userid,
		Order:  order,
		Label:  label,
	}
}

func (g *Group) Save() error {
	return db.Save(g)
}

func (g *Group) UpdateField(field string, change interface{}) error {
	return db.UpdateField(g, field, change)
}

func (g *Group) Delete() error {
	return db.DeleteStruct(g)
}

func (gs Groups) SortByOrder() {
	sort.Slice(gs, func(i, j int) bool {
		return gs[i].Order < gs[j].Order
	})
}

func GroupsAll() Groups {
	var groups Groups
	err := db.All(&groups)
	if err != nil {
		panic(err)
		return groups
	}
	return groups
}

func GroupBy(field string, value interface{}) Group {
	var group Group
	err := db.One(field, value, &group)
	if err != nil {
		panic(err)
		return group
	}
	return group
}

func GroupsBy(field string, value interface{}) Groups {
	var groups Groups
	err := db.Find(field, value, &groups)
	if err != nil {
		panic(err)
		return groups
	}
	return groups
}
