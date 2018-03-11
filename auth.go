package main

import (
	"errors"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Confirm  string `json:"confirm"`
}

func NewCreds(username, password, email, confirm string) *Credentials {
	return &Credentials{
		Username: username,
		Password: password,
		Email:    email,
		Confirm:  confirm,
	}
}

func (creds Credentials) Authenticate() (User, error) {
	var user User
	if creds.Username == "" {
		return user, errors.New("Username is required")
	}
	if creds.Password == "" {
		return user, errors.New("Password is required")
	}
	err := db.QueryRow("SELECT id, username, token FROM users WHERE username=? AND password=?", creds.Username, creds.Password).Scan(&user.Id, &user.Username, &user.Token)
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
	err := db.QueryRow("SELECT username FROM users WHERE username=?", creds.Username).Scan(nil)
	if err == nil {
		return User{}, errors.New("Username already taken")
	}
	var user *User
	user = NewUser(creds.Username, creds.Email, creds.Password)
	err = user.Save()
	if err != nil {
		return *user, err
	}
	return *user, nil
}
