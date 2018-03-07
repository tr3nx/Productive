package main

import (
	"fmt"
	"github.com/asdine/storm"
	"math/rand"
	"time"
)

var (
	db     *storm.DB
	ip     = "127.0.0.1"
	port   = "5000"
	dbpath = "./database.db"
)

func dbConnect(path string, clean bool) *storm.DB {
	if clean {
		deleteIfExists(path)
	}
	db, err := storm.Open(path)
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("[~] starting")

	db = dbConnect(dbpath, false)
	defer db.Close()

	fmt.Println("[@] loading...")
	handleHttp()

	fmt.Println("[!] shutting down!")
}
