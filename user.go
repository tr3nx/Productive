package main

import (
	"fmt"
	"log"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"-"`
	Token    string `json:"token"`
	Created  int64  `json:"-"`
}

type Users []User

var userfields = []string{"id", "username", "password", "email", "token", "created"}

func init() {
	dbRegisterMigration("UsersCreateTable", UsersCreateTable)
	log.Println("[#] User module loading...")
}

func NewUser(username, password, email string) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
		Token:    randomString(16),
		Created:  time.Now().Unix(),
	}
}

func (u *User) Save() error {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO `users`(%v) VALUES(?, ?, ?, ?, ?)", joinFields(userfields[1:])))
	if err != nil {
		return err
	}
	res, err := stmt.Exec(u.Username, u.Password, u.Email, u.Token, u.Created)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = int(id)
	return nil
}

func (u *User) UpdateField(field string, change interface{}) error {
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE `users` %v SET `%v`=`?` WHERE `id`=?", joinFields(userfields), field))
	if err != nil {
		return err
	}
	res, err := stmt.Exec(change, u.Id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error {
	stmt, err := db.Prepare("DELETE FROM `users` WHERE `id`=?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(u.Id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func UsersCreateTable() error {
	stmt, err := db.Prepare("CREATE TABLE `users` (`id` INTEGER PRIMARY KEY AUTO_INCREMENT, `username` VARCHAR(64) NOT NULL, `password` VARCHAR(255) NOT NULL, `email` VARCHAR(64) NOT NULL, `token` VARCHAR(64) NOT NULL, `created` BIGINT NOT NULL)")
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

func UsersDropTable() error {
	stmt, err := db.Prepare("DROP TABLE `users`")
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

func UsersAll() Users {
	var users Users
	rows, err := db.Query(fmt.Sprintf("SELECT %v FROM `users`", joinFields(userfields)))
	if err != nil {
		panic(err)
		return users
	}
	err = rows.Err()
	if err != nil {
		panic(err)
		return users
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Token, &user.Created)
		if err != nil {
			panic(err)
			return users
		}
		users = append(users, user)
	}
	return users
}

func UserBy(field string, value interface{}) User {
	var user User
	err := db.QueryRow(fmt.Sprintf("SELECT %v FROM `users` WHERE `%v`=?", joinFields(userfields), field), value).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Token, &user.Created)
	if err != nil {
		panic(err)
		return user
	}
	return user
}

func UsersBy(field string, value interface{}) Users {
	var users Users
	rows, err := db.Query(fmt.Sprintf("SELECT %v FROM `users` WHERE `%v`=?", joinFields(userfields), field), value)
	if err != nil {
		panic(err)
		return users
	}
	err = rows.Err()
	if err != nil {
		panic(err)
		return users
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Token, &user.Created)
		if err != nil {
			panic(err)
			return users
		}
		users = append(users, user)
	}
	return users
}
