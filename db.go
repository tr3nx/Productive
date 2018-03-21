package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	Migrations   map[string]interface{}
	databaseName = "pro"
)

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

func dbTestData() {
	log.Println("[@] Inserting test data")
	user := NewUser("tr3nx", "me@tr3nx.net", "pass")
	err := user.Save()
	if err != nil {
		log.Println("[?] Test data already loaded?")
		return
	}
	log.Println("[#] Test data loaded")
}

func dbCreateDatabase() error {
	stmt, err := db.Prepare(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", databaseName))
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

func dbDropDatabase() error {
	stmt, err := db.Prepare(fmt.Sprintf("DROP DATABASE IF EXISTS %v", databaseName))
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
