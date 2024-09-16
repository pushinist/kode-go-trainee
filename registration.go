package main

import (
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var existingUser string

		username := r.FormValue("username")
		password := r.FormValue("password")

		err := db.QueryRow(`SELECT username FROM users WHERE username = ?`, username).Scan(&existingUser)
		if existingUser != "" {
			http.Error(w, "User with such name already exists", http.StatusInternalServerError)
			return
		}
		hashedPassword, err := HashPassword(password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		CreateUser(username, hashedPassword)
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("static/signup.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Unable to render page", http.StatusInternalServerError)
		}
	}
}
