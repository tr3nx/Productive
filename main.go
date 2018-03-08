package main

import (
	"github.com/asdine/storm"
	"log"
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
	log.Println("[~] Productive app starting")

	db = dbConnect(dbpath, false)
	defer db.Close()

	log.Println("[@] Loading...")
	handleHttp()

	log.Println("[!] Productive app shutting down!")
}
