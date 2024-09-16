package main

import "net/http"

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/check", spellHandler)
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/add", addHandler)
	mux.HandleFunc("/delete", deleteHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/signup", registerHandler)
	mux.HandleFunc("/logout", logoutHandler)
	http.ListenAndServe(":8080", mux)
}
