package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID   		int64  	`json:"id"`
	Name 		string 	`json:"name"`
	Level 		int64 	`json:"level"`
	Description string 	`json:"description"`
}

type RoleStore struct {
	db *sql.DB
}

func NewRole(db *sql.DB) *RoleStore {
	return &RoleStore{
		db: db,
	}
}

func (r *RoleStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `
		SELECT id, name, level, description
		FROM roles
		WHERE name = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	row := r.db.QueryRowContext(ctx, query, name)
	role := &Role{}
	if err := row.Scan(&role.ID, &role.Name, &role.Level, &role.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return role, nil
}
