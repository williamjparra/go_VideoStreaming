package repository

import (
	"database/sql"
	"errors"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(username, passwordHash string) (*User, error) {
	query := `INSERT INTO users (username, password) VALUES (?,?)`
	result, err := r.DB.Exec(query, username, passwordHash)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Username: username,
	}, nil

}

func (r *UserRepository) GetUserByUsername(username string) (*User, error) {

	query := `SELECT id, username, password FROM users WHERE username = ?`
	row := r.DB.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
