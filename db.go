package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Migrations map[string]interface{}

func init() {
	Migrations = make(map[string]interface{})
}

func dbConnect() *sql.DB {
	db, err := sql.Open("mysql", "homestead:secret@tcp(localhost:33060)/pro")
	if err != nil {
		panic(err)
	}
	return db
}

func dbRegisterMigration(name string, qryFunc func() error) {
	Migrations[name] = qryFunc
}

func dbMigrate() {
	log.Println("[@] Migrations running")
	var err error
	for name, qryFn := range Migrations {
		err = qryFn.(func() error)()
		if err != nil {
			log.Println(fmt.Sprintf("[!] \"%v\" migration failed", name))
			continue
		}
		log.Println(fmt.Sprintf("[+] \"%v\" migrated", name))
	}
	log.Println("[#] Migrations completed")
}
