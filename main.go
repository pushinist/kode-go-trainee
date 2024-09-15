package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// add notes ✅
// read notes ✅
// validate notes (ya.speller)
// authentication and authorization
// special access
// registration
// logging
// docker

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	notes, err := GetNotes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, notes)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")

		CreateNote(title, content)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	http.ListenAndServe(":8080", nil)
}
