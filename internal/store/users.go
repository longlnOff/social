package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID 			int64 `json:"id"`
	Username	string `json:"username"`
	Email		string `json:"email"`
	Password	string `json:"-"`
	CreatedAt	string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (u *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1,$2,$3)
		RETURNING id, created_at
	`
	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return u.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
}

func (u *UserStore) GetByUserID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := u.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrNotFound
			default:
				return nil, err
		}
	}

	return &user, nil
}


