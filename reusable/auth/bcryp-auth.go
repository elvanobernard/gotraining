package auth

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func TryLogin(db *sql.DB, username, password string) (bool, error) {
	var hash string
	row := db.QueryRow("SELECT hash_password FROM users WHERE user_name = ?", username)
	if err := row.Scan(hash); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("no user found")
		}

		return false, fmt.Errorf("error during login")
	}

	return CheckPasswordHash(password, hash), nil
}

func NewUser(db *sql.DB, username, password string) (int64, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("error during creating hash password")
	}

	result, err := db.Exec("INSERT INTO users (user_name, hash_password) VALUES (?, ?)", username, hash)

	if err != nil {
		return 0, fmt.Errorf("error during inserting to database")
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("error during retrieving ID")
	}

	return id, nil
}
