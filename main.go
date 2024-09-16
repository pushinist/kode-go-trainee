package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	Username string
	Notes    []Note
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	var username string
	var userID int
	if err == nil {
		username = sessions[cookie.Value]
		row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)
		err = row.Scan(&userID)
	}

	rows, err := db.Query("SELECT id, title, content FROM notes WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var notes []Note
	for rows.Next() {
		var note Note
		if err = rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			http.Error(w, "Unable to scan note", http.StatusInternalServerError)
			return
		}
		note.UserID = userID
		notes = append(notes, note)
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err = tmpl.Execute(w, PageData{
		Username: username,
		Notes:    notes,
	})
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		username := sessions[cookie.Value]
		row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)
		var userID int
		err = row.Scan(&userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "error scanning row", http.StatusInternalServerError)
			return
		}
		CreateNote(title, content, userID)
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	DeleteNote(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	initDB()
	defer db.Close()
	http.HandleFunc("/check", spellHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", registerHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.ListenAndServe(":8080", nil)
}
