package models

import (
	"database/sql"
	"time"
)

type Note struct {
	ID       int
	Title    string
	Content  string
	Created  time.Time
	Expirese time.Time
}

type NoteModel struct {
	DB *sql.DB
}

func (m *NoteModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func (m *NoteModel) Get(id int) (*Note, error) {
	return nil, nil
}

func (m *NoteModel) Latest() ([]*Note, error) {
	return nil, nil
}
