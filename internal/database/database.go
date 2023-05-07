package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/EgorKo25/GophKeeper/internal/storage"
)

// ManagerDB structure for managing database
type ManagerDB struct {
	Db *sql.DB
}

// NewManagerDB constructor
func NewManagerDB(address string) (*ManagerDB, error) {

	ctx := context.Background()

	db, err := sql.Open("pgx", address)
	if err != nil {
		return nil, err
	}

	err = createTable(ctx, db)
	if err != nil {
		return nil, err
	}

	return &ManagerDB{
		Db: db,
	}, nil
}

// createTable creates tables for the database operation
func createTable(ctx context.Context, db *sql.DB) error {

	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := `CREATE IF NOT EXISTS 
    TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := db.ExecContext(childCtx, query)
	if err != nil {
		return err
	}

	return nil
}

// Ping testing connection to database
func (m *ManagerDB) Ping() bool {
	err := m.Db.Ping()
	return err == nil
}

// AddUser adds new user to the database
func (m *ManagerDB) AddUser(ctx context.Context, user *storage.User) error {

	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := `INSERT INTO users (username, password, email, created_at, updated_at) 
							VALUES  ($1, $2, $3, $4, $5)`
	_, err := m.Db.ExecContext(childCtx, query,
		user.Login, user.Passwd, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
