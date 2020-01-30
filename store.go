package main

import "database/sql"

type SimpleStore struct {
	db *sql.DB
}

func NewSimpleStore(db *sql.DB) *SimpleStore {
	return &SimpleStore{
		db: db,
	}
}

type User struct {
	Username string
	PasswordHash []byte
}

func (s *SimpleStore) GetUserByUsername(username string) (*User, error) {
	//
	var user User
	//
	row := s.db.QueryRow(`SELECT username, password_hash FROM users WHERE username = $1`, username)
	//
	if err := row.Scan(&user.Username, &user.PasswordHash); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SimpleStore) CreateSession() (string, error) {
	//
	var token string
	row := s.db.QueryRow(`INSERT INTO sessions DEFAULT VALUES RETURNING token;`)
	if err := row.Scan(&token); err != nil {
		return "", err
	}

	return token, nil
}
