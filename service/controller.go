package service

import (
	"context"

	pb "github.com/meateam/permission-service/proto"
)

// Controller is an interface for the business logic of the permission.Service which uses a Store.
type Controller interface {
	CreatePermission(
		ctx context.Context,
		fileID string,
		userID string,
		role pb.Role,
		creator string,
		override bool) (Permission, error)
	DeletePermission(ctx context.Context, fileID string, userID string) (Permission, error)
	GetFilePermissions(ctx context.Context, fileID string) ([]*pb.GetFilePermissionsResponse_UserRole, error)
	GetByFileAndUser(ctx context.Context, fileID string, userID string) (Permission, error)
	GetUserPermissions(ctx context.Context, userID string) ([]*pb.GetUserPermissionsResponse_FileRole, error)
	DeleteFilePermissions(ctx context.Context, fileID string) ([]*pb.PermissionObject, error)
	HealthCheck(ctx context.Context) (bool, error)
}
