package service

import (
	"context"
)

// Controller is an interface for the business logic of the permission.Service which uses a Store.
type Controller interface {
	CreatePermission(ctx context.Context, fileID string, userID string) (string, error)
	GetPermission(ctx context.Context, fileID string, userID string) (Permission, error)
	DeletePermission(ctx context.Context, fileID string, userID string) (Permission, error)
	HealthCheck(ctx context.Context) (bool, error)
}
