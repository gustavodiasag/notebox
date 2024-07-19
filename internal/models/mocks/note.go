package mocks

import (
	"time"

	"github.com/gustavodiasag/notebox/internal/models"
)

var mockNote = &models.Note{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond",
	Created: time.Now(),
	Expires: time.Now(),
}

type NoteModel struct{}

func (m *NoteModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *NoteModel) Get(id int) (*models.Note, error) {
	switch id {
	case 1:
		return mockNote, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *NoteModel) Latest() ([]*models.Note, error) {
	return []*models.Note{mockNote}, nil
}
