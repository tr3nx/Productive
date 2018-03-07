package main

type User struct {
	Id       int    `json:"id" storm:"id,increment"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: HashPassword(password),
	}
}
