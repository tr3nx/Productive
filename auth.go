package main

import (
	"github.com/asdine/storm/q"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewCreds(username, password string) *Credentials {
	return &Credentials{
		Username: username,
		Password: password,
	}
}

func (c Credentials) Authenticate() (User, error) {
	var user User
	err := db.Select(q.Eq("Username", c.Username), q.Eq("Password", HashPassword(c.Password))).First(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
