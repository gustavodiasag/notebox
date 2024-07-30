package mocks

import (
	"time"

	"github.com/gustavodiasag/notebox/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "foo@mail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	if id != 1 {
		return nil, models.ErrNoRecord
	}

	u := &models.User{
		ID:      1,
		Name:    "Alice",
		Email:   "alice@example.com",
		Created: time.Now(),
	}
	return u, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pass" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) UpdatePassword(id int, current, new string) error {
	if id != 1 {
		return models.ErrNoRecord
	}
	if current != "pass" {
		return models.ErrInvalidCredentials
	}

	return nil
}
