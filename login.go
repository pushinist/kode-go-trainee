package main

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

var sessions = map[string]string{}

func checkUser(username, password string) (string, error) {
	var storedHash string
	err := db.QueryRow(`SELECT password FROM users WHERE username = ?`, username).Scan(&storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "user not found", fmt.Errorf("user not found")
		}
		return "something went wrong", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return "wrong password", fmt.Errorf("wrong password")
	}
	return "", nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if msg, err := checkUser(username, password); err != nil {
			http.Error(w, msg, http.StatusInternalServerError)
		}
		randomSessionToken, err := generateToken()
		if err != nil {
			http.Error(w, "Error generating session token", http.StatusInternalServerError)
		}
		sessions[randomSessionToken] = username
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    randomSessionToken,
			HttpOnly: true,
			Path:     "/"})
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("static/login.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Unable to render page", http.StatusInternalServerError)
		}
	}
	return

}
