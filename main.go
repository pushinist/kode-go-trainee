package main

type PageData struct {
	Username string
	Notes    []Note
}

func main() {
	initDB()
	defer db.Close()
	startServer()
}
