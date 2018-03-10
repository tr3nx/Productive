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
	if creds.Username == "" {
		return User{}, errors.New("Username is required")
	}
	if creds.Password == "" {
		return User{}, errors.New("Password is required")
	}
	var user User
	err := db.Select(q.Eq("Username", creds.Username), q.Eq("Password", creds.Password)).First(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (creds Credentials) Signup() (User, error) {
	if creds.Username == "" {
		return User{}, errors.New("Username is required")
	}
	if creds.Email == "" {
		return User{}, errors.New("Email is required")
	}
	if creds.Password == "" {
		return User{}, errors.New("Password is required")
	}
	if creds.Password != creds.Confirm {
		return User{}, errors.New("Passwords do not match")
	}

	err := db.One("Username", creds.Username, nil)
	if err == nil {
		return User{}, errors.New("Username already taken")
	}
	user := NewUser(creds.Username, creds.Email, creds.Password)
	err = user.Save()
	if err != nil {
		return User{}, err
	}
	return *user, nil
}
