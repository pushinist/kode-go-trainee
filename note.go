package main

import "log"

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

func GetNotes() ([]Note, error) {
	rows, err := db.Query("SELECT id, title, content FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var note Note
		err = rows.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

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
