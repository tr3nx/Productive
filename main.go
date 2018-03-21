package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"
)

var (
	db    *sql.DB
	ip    = "127.0.0.1"
	port  = "5000"
	clean = true
)

func init() {
	log.Println("=== Productive app starting")
	rand.Seed(time.Now().UnixNano())
	db = dbConnect()
	log.Println("[+] Connected to database.")
	log.Println("[$] Initializing modules...")
}

func main() {
	defer db.Close()

	if clean {
		dbDropDatabase()
		dbCreateDatabase()
		dbMigrate()
		dbTestData()
	}

	handleHttp()

	log.Println("[!] Productive app shutting down!")
}
