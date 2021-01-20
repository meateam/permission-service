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
		override bool,
		appID string) (Permission, error)
	DeletePermission(ctx context.Context, fileID string, userID string) (Permission, error)
	GetFilePermissions(ctx context.Context, fileID string) ([]*pb.GetFilePermissionsResponse_UserRole, error)
	GetByFileAndUser(ctx context.Context, fileID string, userID string) (Permission, error)
	GetPermissionByMongoID(ctx context.Context, mongoID string) (Permission, error)
	GetUserPermissions(ctx context.Context, userID string, pageNum int64, pageSize int64, isShared bool, appID string) (*pb.GetUserPermissionsResponse, error)
	DeleteFilePermissions(ctx context.Context, fileID string) ([]*pb.PermissionObject, error)
	HealthCheck(ctx context.Context) (bool, error)
}
