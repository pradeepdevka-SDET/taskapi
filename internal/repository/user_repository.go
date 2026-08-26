package repository

import (
	"database/sql"
	"taskapi/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(username, passwordHash string) (int, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO users (username, password_hash) VALUES ($1,$2) RETURNING id`,
		username, passwordHash,
	).Scan(&id)
	return id, err
}

func (r *UserRepository) GetByUsername(username string) (model.User, error) {
	var u model.User
	err := r.db.QueryRow(
		`SELECT id,username,password_hash FROM users WHERE username=$1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash)
	return u, err
}
