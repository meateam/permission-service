package service

import (
	"context"
)

// Store is an interface for handling the storing of permissions.
type Store interface {
	Create(ctx context.Context, permission Permission) (Permission, error)
	Get(ctx context.Context, filter interface{}) (Permission, error)
	GetAll(ctx context.Context, filter interface{}) ([]Permission, error)
	Delete(ctx context.Context, filter interface{}) (Permission, error)
	HealthCheck(ctx context.Context) (bool, error)
}
