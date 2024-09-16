package main

import (
	"crypto/rand"
	"encoding/hex"
)

func CreateUser(username, password string) {
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err := db.Exec(query, username, password)
	checkErr(err)
}

func generateToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(tokenBytes)
	return token, nil
}
