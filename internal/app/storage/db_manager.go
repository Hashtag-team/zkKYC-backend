package storage

import (
	"context"
	"database/sql"
	"fmt"
)

var (
	CreateUsersTableRaw = `
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(1) NOT NULL CHECK (gender IN ('M', 'F')),
    eth_address VARCHAR(42) NOT NULL UNIQUE,
    did VARCHAR(255) UNIQUE CHECK (did ~ '^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$')
);
CREATE UNIQUE INDEX idx_full_name ON users (full_name);
`

	CreateUser = `
INSERT INTO users (full_name, birth_date, gender, eth_address, did)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id;
`
	GetByDid = `
        SELECT id, full_name, birth_date, gender, eth_address, did
        FROM users
        WHERE eth_address = $1;
    `
)

// Connection to database
type DBManager struct {
	DB *sql.DB
}

// Function to create new manager
func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{DB: db}
}

// Create table
func (m *DBManager) CreateTable(ctx context.Context) error {
	_, err := m.DB.ExecContext(ctx, CreateUsersTableRaw)
	return err
}

// CreateUser
func (m *DBManager) CreateUser(ctx context.Context, user *User) error {
	return m.DB.QueryRowContext(ctx, CreateUser,
		user.FullName,
		user.BirthDate,
		user.Gender,
		user.EthAddress,
		user.DID,
	).Scan(&user.ID)
}

// GetUserByAddress returns user with specific did
func (m *DBManager) GetUserByAddress(ctx context.Context, eth string) (*User, error) {

	var user User
	err := m.DB.QueryRowContext(ctx, GetByDid, eth).Scan(
		&user.ID,
		&user.FullName,
		&user.BirthDate,
		&user.Gender,
		&user.EthAddress,
		&user.DID,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by did error: %w", err)
	}

	return &user, nil
}
