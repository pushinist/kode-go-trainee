package main

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Note struct {
	ID      int
	Title   string
	Content string
	UserID  int
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreateNote(title, content string, userIDInt int) {
	query := `INSERT INTO notes (title, content, user_id) VALUES (?, ?, ?)`
	_, err := db.Exec(query, title, content, userIDInt)
	checkErr(err)
}

func GetNotes(userID int) ([]Note, error) {
	rows, err := db.Query("SELECT id, title, content FROM notes WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []Note
	for rows.Next() {
		var note Note
		if err = rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			return nil, err
		}
		note.UserID = userID
		notes = append(notes, note)
	}
	return notes, nil
}

// На будущее, для возможности редактировать заметки
func UpdateNotes(id int, title, content string) {
	query := `UPDATE notes SET title = ?, content = ? WHERE id = ?`
	_, err := db.Exec(query, title, content, id)
	checkErr(err)
}

func DeleteNote(id int) {
	query := `DELETE FROM notes WHERE id = ?`
	_, err := db.Exec(query, id)
	checkErr(err)
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
	notes, err := GetNotes(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
