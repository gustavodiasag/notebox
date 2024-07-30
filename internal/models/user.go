package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Get(id int) (*User, error)
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	UpdatePassword(id int, current, new string) error
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	// The second parameter passed to `GenerateFromPassword` indicates the number of bcrypt
	// iterations used to generate the password hash.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
    INSERT INTO user (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())
  `
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// Check specifically if the error generated is related to the unique email constraint.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "user_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Get(id int) (*User, error) {
	var user User

	stmt := `SELECT id, name, email, created FROM user WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM user WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT true FROM user WHERE id = ?)`

	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}

func (m *UserModel) UpdatePassword(id int, current, new string) error {
	var currentHashedPassword []byte

	stmt := `SELECT hashed_password FROM user WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&currentHashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(current))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(new), 12)
	if err != nil {
		return err
	}

	stmt = `UPDATE user SET hashed_password = ? WHERE id = ?`

	_, err = m.DB.Exec(stmt, newHashedPassword, id)
	return err
}
