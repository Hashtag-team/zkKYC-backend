package storage

import "context"

type Repository interface {
	Add(ctx context.Context, user *User) error
	Get(ctx context.Context, address string) (interface{}, bool)
}
