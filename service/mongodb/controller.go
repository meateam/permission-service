package mongodb

import (
	"context"
	"fmt"

	pb "github.com/meateam/permission-service/proto"
	"github.com/meateam/permission-service/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// IsPermitted returns the permission in store that matches fileID and userID.
func (c Controller) IsPermitted(ctx context.Context, fileID string, userID string, role pb.Role) (bool, error) {
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
	if err != nil {
		return false, err
	}

	return isSubRole(permission.GetRole(), role), nil
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
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// HealthCheck runs store's healthcheck and returns true if healthy, otherwise returns false
// and any error if occured.
func (c Controller) HealthCheck(ctx context.Context) (bool, error) {
	return c.store.HealthCheck(ctx)
}

// GetFilePermissions returns a slice of UserRoles,
// otherwise returns nil and any error if occured.
func (c Controller) GetFilePermissions(ctx context.Context, fileID string) ([]*pb.GetFilePermissionsResponse_UserRole, error) {
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

func isSubRole(role pb.Role, wanted pb.Role) bool {
	if wanted == pb.Role_NONE {
		return false
	}

	switch role {
	case pb.Role_NONE:
		return false
	case pb.Role_OWNER:
		return true
	case pb.Role_WRITE:
		return wanted == pb.Role_WRITE || wanted == pb.Role_READ
	case pb.Role_READ:
		return wanted == pb.Role_READ
	default:
		return false
	}
}
