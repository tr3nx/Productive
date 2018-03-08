package main

import (
	"errors"
	"github.com/asdine/storm/q"
)

type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func NewCreds(username, email, password, confirm string) *Credentials {
	return &Credentials{
		Username: username,
		Email:    email,
		Password: password,
		Confirm:  confirm,
	}
}

func (creds Credentials) Authenticate() (User, error) {
	var user User
	err := db.Select(q.Eq("Username", creds.Username), q.Eq("Password", creds.Password)).First(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (creds Credentials) Signup() (User, error) {
	err := db.One("Username", creds.Username, nil)
	if err == nil {
		return User{}, errors.New("Username already taken")
	}
	if creds.Password != creds.Confirm {
		return User{}, errors.New("Passwords do not match")
	}
	user := NewUser(creds.Username, creds.Email, creds.Password)
	err = user.Save()
	if err != nil {
		return User{}, err
	}
	return *user, nil
}
