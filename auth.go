package main

import (
	"github.com/asdine/storm/q"
)

type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
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

func (c Credentials) Signup() (User, error) {
	var user User
	err := db.Select(q.Eq("Username", c.Username), q.Eq("Password", HashPassword(c.Password))).First(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
