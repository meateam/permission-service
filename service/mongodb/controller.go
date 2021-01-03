package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/meateam/permission-service/proto"
	"github.com/meateam/permission-service/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Controller is the permissions service business logic implementation using MongoStore.
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
func (c Controller) CreatePermission(
	ctx context.Context,
	fileID string,
	userID string,
	role pb.Role,
	creator string,
	override bool,
	appID string) (service.Permission, error) {
	permission := &BSON{FileID: fileID, UserID: userID, Role: role, Creator: creator, AppID: appID}
	createdPermission, err := c.store.Create(ctx, permission, override)
	if err != nil {
		return nil, fmt.Errorf("failed creating permission: %v", err)
	}

	return createdPermission, nil
}

// GetByFileAndUser retrieves the permissoin that matches fileID and userID, and any error if occurred.
func (c Controller) GetByFileAndUser(
	ctx context.Context,
	fileID string,
	userID string) (service.Permission, error) {
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
		return nil, status.Error(codes.NotFound, "permission not found")
	}

	return permission, nil
}

// GetPermissionByMongoID retrieves the permissoin by the recived mongo id.
func (c Controller) GetPermissionByMongoID(
	ctx context.Context,
	mongoID string) (service.Permission, error) {
	filter := bson.D{
		bson.E{
			Key:   MongoObjectIDField,
			Value: mongoID,
		},
	}

	permission, err := c.store.Get(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "permission not found")
	}

	return permission, nil
}

// DeletePermission deletes the permission in store that matches fileID and userID
// and returns the deleted permission.
func (c Controller) DeletePermission(
	ctx context.Context,
	fileID string,
	userID string,
) (service.Permission, error) {
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
		return nil, status.Error(codes.NotFound, "permission not found")
	}

	return permission, nil
}

// HealthCheck runs store's healthcheck and returns true if healthy, otherwise returns false
// and any error if occurred.
func (c Controller) HealthCheck(ctx context.Context) (bool, error) {
	return c.store.HealthCheck(ctx)
}

// GetFilePermissions returns a slice of UserRole,
// otherwise returns nil and any error if occurred.
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
			UserID:  permission.GetUserID(),
			Role:    permission.GetRole(),
			Creator: permission.GetCreator(),
		})
	}
	return returnedPermissions, nil
}

// GetUserPermissions returns a slice of FileRole,
// otherwise returns nil and any error if occurred.
func (c Controller) GetUserPermissions(
	ctx context.Context,
	userID string, pageNum int64, pageSize int64, isShared bool, appID string) (*pb.GetUserPermissionsResponse, error) {

	// Check if one is negative and the other is not
	if pageNum < 0 || pageSize < 0 {
		return nil, fmt.Errorf("pageNum %d and pageSize %d must both be non-negative", pageNum, pageSize)
	}

	var filter bson.D

	filter = append(filter, bson.E{
		Key:   PermissionBSONUserIDField,
		Value: userID,
	})

	if isShared {
		filter = append(filter, bson.E{
			Key:   PermissionBSONCreatorField,
			Value: bson.M{"$ne": userID},
		})
	}

	if appID != "" {
		filter = append(filter, bson.E{
			Key:   PermissionBSONAppIDField,
			Value: appID,
		})
	}

	sort := bson.D{
		bson.E{
			Key:   MongoObjectIDField,
			Value: 1,
		},
	}

	// Get permissions by page, sorted by mongoID
	pageRes, err := c.store.GetUserPermissionsByPage(ctx, pageNum, pageSize, sort, filter)
	if err != nil {
		return nil, err
	}

	// Go over the page of permissions received and reformat them
	filePermissions := c.reformatFilePermissions(pageRes.permissions)

	return &pb.GetUserPermissionsResponse{Permissions: filePermissions, ItemCount: pageRes.itemCount, PageNum: pageRes.pageNum}, nil

}

// reformatFilePermissions receives an array of service.Permission and
// asynchronously returns them as an array of *pb.GetUserPermissionsResponse_FileRole,
// while keeping the order in which they were received

func (c Controller) reformatFilePermissions(permissions []service.Permission) []*pb.GetUserPermissionsResponse_FileRole {
	start := time.Now().UnixNano()
	filePermissions := make([]*pb.GetUserPermissionsResponse_FileRole, len(permissions), len(permissions))

	var wg sync.WaitGroup
	wg.Add(len(permissions))

	// for i := 0; i < len(permissions); i++ {
	// 	permission := permissions[i]
	// 	go func(i int, permission service.Permission) {
	// 		defer wg.Done()
	// 		time.Sleep(100 * time.Nanosecond)
	// 		filePermissions[i] = &pb.GetUserPermissionsResponse_FileRole{
	// 			FileID:  permission.GetFileID(),
	// 			Role:    permission.GetRole(),
	// 			Creator: permission.GetCreator(),
	// 		}
	// 	}(i, permission)
	// }
	// wg.Wait()

	// 10 Nano sleep:
	// Finished for loop 885832 Async
	// Finished for loop 121066 Sync

	// 100 Nano sleep:
	// Finished for loop 264884 Async
	// Finished for loop 579209 Sync

	// 13953001
	// 10,322,789
	// 1,117,611
	// 159,598
	for _, permission := range permissions {
		// time.Sleep(100 * time.Nanosecond)
		// start2 := time.Now().UnixNano()
		filePermissions = append(filePermissions, &pb.GetUserPermissionsResponse_FileRole{
			FileID:  permission.GetFileID(),
			Role:    permission.GetRole(),
			Creator: permission.GetCreator(),
		})
		// end2 := time.Now().UnixNano()
		// fmt.Printf("Finished for loop %v\n", end2-start2)
	}
	end := time.Now().UnixNano()

	fmt.Printf("Finished for loop %v\n", end-start)

	return filePermissions
}

// func (c Controller) reformatFilePermissions(permissions []service.Permission) []*pb.GetUserPermissionsResponse_FileRole {
// 	start := time.Now().UnixNano()
// 	filePermissions := make([]*pb.GetUserPermissionsResponse_FileRole, len(permissions), len(permissions))

// 	for _, permission := range permissions {
// 		filePermissions = append(filePermissions, &pb.GetUserPermissionsResponse_FileRole{
// 			FileID:  permission.GetFileID(),
// 			Role:    permission.GetRole(),
// 			Creator: permission.GetCreator(),
// 		})
// 	}
// 	end := time.Now().UnixNano()

// 	fmt.Printf("Finished for loop %v\n", end-start)

// 	return filePermissions
// }

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
		permissionID, err := primitive.ObjectIDFromHex(permission.GetID())
		if err != nil {
			return nil, err
		}

		permissionFilter := bson.D{
			bson.E{
				Key:   MongoObjectIDField,
				Value: permissionID,
			},
		}

		deletedPermission, err := c.store.Delete(ctx, permissionFilter)
		if err != nil {
			return nil, err
		}

		protoDeletedPermission := &pb.PermissionObject{
			Id:      deletedPermission.GetID(),
			FileID:  deletedPermission.GetFileID(),
			UserID:  deletedPermission.GetUserID(),
			Role:    deletedPermission.GetRole(),
			Creator: deletedPermission.GetCreator(),
		}
		deletedPermissions = append(deletedPermissions, protoDeletedPermission)
	}

	return deletedPermissions, nil
}
