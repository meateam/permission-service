package mongodb

import (
	"context"
	"fmt"

	pb "github.com/meateam/permission-service/proto"
	"github.com/meateam/permission-service/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Controller is the permisison service business logic implementation using MongoStore.
type Controller struct {
	store MongoStore
}

// NewMongoController returns a new controller.
func NewMongoController(db *mongo.Database) (Controller, error) {
	store, err := newMongoStore(db)
	if err != nil {
		return Controller{}, err
	}

	return Controller{store: store}, nil
}

// CreatePermission creates a Permission in store and returns its unique ID.
func (c Controller) CreatePermission(ctx context.Context, fileID string, userID string, role pb.Role) (service.Permission, error) {
	// Create the root permission.
	permission := &BSON{FileID: fileID, UserID: userID, Role: role}
	createdPermission, err := c.store.Create(ctx, permission)
	if err != nil {
		return nil, fmt.Errorf("failed creating permission: %v", err)
	}

	return createdPermission, nil
}

// GetByFileAndUser retrieves the permissoin that matches fileID and userID, and any error if occured.
func (c Controller) GetByFileAndUser(ctx context.Context, fileID string, userID string) (service.Permission, error) {
	filter := bson.D{
		bson.E{
			Key:   PermissionBSONFileIDField,
			Value: fileID,
		},
		bson.E{
			Key:   PermissionBSONUserIDField,
			Value: userID,
		},
	}

	permission, err := c.store.Get(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.Unimplemented, "permission not found")
	}

	return permission, nil
}

// DeletePermission deletes the permission in store that matches fileID and userID
// and returns the deleted permission.
func (c Controller) DeletePermission(ctx context.Context, fileID string, userID string) (service.Permission, error) {
	filter := bson.D{
		bson.E{
			Key:   PermissionBSONFileIDField,
			Value: fileID,
		},
		bson.E{
			Key:   PermissionBSONUserIDField,
			Value: userID,
		},
	}

	permission, err := c.store.Delete(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.Unimplemented, "permission not found")
	}

	return permission, nil
}

// HealthCheck runs store's healthcheck and returns true if healthy, otherwise returns false
// and any error if occured.
func (c Controller) HealthCheck(ctx context.Context) (bool, error) {
	return c.store.HealthCheck(ctx)
}

// GetFilePermissions returns a slice of UserRole,
// otherwise returns nil and any error if occured.
func (c Controller) GetFilePermissions(ctx context.Context,
	fileID string) ([]*pb.GetFilePermissionsResponse_UserRole, error) {
	filter := bson.D{
		bson.E{
			Key:   PermissionBSONFileIDField,
			Value: fileID,
		},
	}

	filePermissions, err := c.store.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	returnedPermissions := make([]*pb.GetFilePermissionsResponse_UserRole, 0, len(filePermissions))
	for _, permission := range filePermissions {
		returnedPermissions = append(returnedPermissions, &pb.GetFilePermissionsResponse_UserRole{
			UserID: permission.GetUserID(),
			Role:   permission.GetRole(),
		})
	}
	return returnedPermissions, nil
}

// GetUserPermissions returns a slice of FileRole,
// otherwise returns nil and any error if occured.
func (c Controller) GetUserPermissions(ctx context.Context,
	userID string, isOwner bool) ([]*pb.GetUserPermissionsResponse_FileRole, error) {
	filter := bson.D{
		bson.E{
			Key:   PermissionBSONUserIDField,
			Value: userID,
		},
	}

	if isOwner {
		filter = append(filter, bson.E{
			Key:   PermissionBSONRoleField,
			Value: pb.Role_OWNER,
		})
	} else {
		filter = append(filter,
			bson.E{
				Key: PermissionBSONRoleField,
				Value: bson.M{
					"$nin": bson.A{pb.Role_OWNER, pb.Role_NONE},
				},
			},
		)
	}

	permissions, err := c.store.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	filePermissions := make([]*pb.GetUserPermissionsResponse_FileRole, 0, len(permissions))
	for _, permission := range permissions {
		filePermissions = append(filePermissions, &pb.GetUserPermissionsResponse_FileRole{
			FileID: permission.GetFileID(),
			Role:   permission.GetRole(),
		})
	}

	return filePermissions, nil
}

// DeleteFilePermissions deletes all permissions that exist for fileID and
// returns a slice of Permissions that were deleted.
func (c Controller) DeleteFilePermissions(ctx context.Context,
	fileID string) ([]*pb.PermissionObject, error) {
	filePermissionsFilter := bson.D{
		bson.E{
			Key:   PermissionBSONFileIDField,
			Value: fileID,
		},
	}
	permissions, err := c.store.GetAll(ctx, filePermissionsFilter)
	if err != nil {
		return nil, err
	}

	deletedPermissions := make([]*pb.PermissionObject, 0, len(permissions))
	for _, permission := range permissions {
		permissionFilter := bson.D{
			bson.E{
				Key:   MongoObjectIDField,
				Value: permission.GetFileID(),
			},
		}

		deletedPermission, err := c.store.Delete(ctx, permissionFilter)
		if err != nil {
			return nil, err
		}

		protoDeletedPermission := &pb.PermissionObject{
			Id:     deletedPermission.GetID(),
			FileID: deletedPermission.GetFileID(),
			UserID: deletedPermission.GetUserID(),
			Role:   deletedPermission.GetRole(),
		}
		deletedPermissions = append(deletedPermissions, protoDeletedPermission)
	}

	return deletedPermissions, nil
}
