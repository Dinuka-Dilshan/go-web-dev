package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RoleStore struct {
	db *pgxpool.Pool
}

func (roleStore *RoleStore) GetByName(ctx context.Context, name string) (*Role, error) {

	query := `SELECT id, name, level, description FROM roles WHERE name = $1`
	var role Role

	err := roleStore.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &role.Level, &role.Description)
	if err != nil {
		return nil, err
	}

	return &role, nil
}
