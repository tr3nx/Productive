package main

type User struct {
	Id       int    `json:"id" storm:"id,increment"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Users []User

func NewUser(username, email, password string) *User {
	return &User{
		Username: username,
		Email:    email,
		Password: password,
		Token:    randomString(16),
	}
}

func (u *User) Save() error {
	return db.Save(u)
}

func (u *User) UpdateField(field string, change interface{}) error {
	return db.UpdateField(u, field, change)
}

func (u *User) Delete() error {
	return db.DeleteStruct(u)
}

func UsersAll() Users {
	var users Users
	err := db.All(&users)
	if err != nil {
		panic(err)
		return users
	}
	return users
}

func UserBy(field string, value interface{}) User {
	var user User
	err := db.One(field, value, &user)
	if err != nil {
		panic(err)
		return user
	}
	return user
}

func UsersBy(field string, value interface{}) Users {
	var users Users
	err := db.Find(field, value, &users)
	if err != nil {
		panic(err)
		return users
	}
	return users
}
