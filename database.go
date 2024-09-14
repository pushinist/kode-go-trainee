package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error

	db, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal(err)
	}
	createTable := `
CREATE TABLE IF NOT EXISTS notes (
id INTEGER PRIMARY KEY AUTOINCREMENT,
title TEXT,
content TEXT);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}
