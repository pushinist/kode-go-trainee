package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error

	db, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	createNotesTable := `
CREATE TABLE IF NOT EXISTS notes (
id INTEGER PRIMARY KEY AUTOINCREMENT,
title TEXT,
content TEXT,
user_id INTEGER,
FOREIGN KEY (user_id) REFERENCES users(id));`

	createUsersTable := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    password TEXT);
`

	_, err = db.Exec(createNotesTable)
	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}
}
