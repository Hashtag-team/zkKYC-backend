package storage

import (
	"context"
	"database/sql"
	"time"
)

// URL struct
type User struct {
	ID         int       `db:"id" json:"id"`
	FullName   string    `db:"full_name" json:"full_name"`
	BirthDate  time.Time `db:"birth_date" json:"birth_date"`
	Gender     string    `db:"gender" json:"gender"`
	EthAddress string    `db:"eth_address" json:"eth_address"`
	DID        string    `db:"did" json:"did"`
}

// DBStorage stores pointer to DBManager
type DBStorage struct {
	dbManager *DBManager
}

// Create new DB Storage
func NewDBStorage(db *sql.DB) Repository {
	dbs := DBStorage{
		dbManager: NewDBManager(db),
	}

	dbs.dbManager.CreateTable(context.Background())

	return &dbs
}

// Add value
func (s *DBStorage) Add(ctx context.Context, user *User) error {
	return s.dbManager.CreateUser(ctx, user)
}

// Get value
func (s *DBStorage) Get(ctx context.Context, address string) (interface{}, bool) {
	v, err := s.dbManager.GetUserByAddress(ctx, address)
	if err != nil {
		return "", false
	}
	return v, true
}
