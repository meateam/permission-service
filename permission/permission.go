package permission

import (
	"context"
	"fmt"
	"time"

	pb "github.com/meateam/permission-service/proto"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// PermissionCollectionName is the name of the permissions collection.
	PermissionCollectionName = "permissions"
)

// Service is a structure used for handling Permission Service grpc requests.
type Service struct {
	db     *mongo.Database
	logger *logrus.Logger
}

// Permission is the structure that represents a permission as it's stored.
type Permission struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FileID    string             `json:"fileID,omitempty" bson:"fileID,omitempty"`
	UserID    string             `json:"userID,omitempty" bson:"userID,omitempty"`
	Inherited primitive.ObjectID `json:"inherited,omitempty" bson:"inherited,omitempty"`
}

// Store is an interface for handling the storing of permissions.
type Store interface {
	Create(ctx context.Context, permission Permission) (string, error)
	CreateMany(ctx context.Context, permissions []Permission) ([]string, error)
	Get(ctx context.Context, permissionID string) (Permission, error)
	GetByFileAndUser(ctx context.Context, fileID string, userID string) (Permission, error)
	Delete(ctx context.Context, permissionID string) (Permission, error)
}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s *Service) HealthCheck(mongoClientPingTimeout time.Duration) bool {
	if s == nil {
		panic("s == nil")
	}

	timeoutCtx, cancel := context.WithTimeout(context.TODO(), mongoClientPingTimeout)
	defer cancel()
	err := s.db.Client().Ping(timeoutCtx, readpref.Primary())
	if err != nil {
		s.logger.Errorf("%v", err)
		return false
	}

	return true
}

// NewService creates a Service and returns it.
func NewService(db *mongo.Database, logger *logrus.Logger) *Service {
	return &Service{db: db, logger: logger}
}

// CreatePermission is the request handler for creating a permission of a file to user.
func (s Service) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.CreatePermissionResponse, error) {
	rootFileID := req.GetFileID()
	userID := req.GetUserID()

	permission, err := s.GetByFileAndUser(ctx, rootFileID, userID)
	if err != nil {
		return nil, err
	}

	if permission.ID.IsZero() {
		return nil, fmt.Errorf("user %s is already permitted to file %s", userID, rootFileID)
	}

	// Create the root permission.
	rootPermission := Permission{FileID: rootFileID, UserID: userID}
	rootPermissionID, err := s.Create(ctx, rootPermission)
	if err != nil {
		s.logger.Errorf("failed creating root permission: %v", err)
		return nil, err
	}

	inheritors := make([]Permission, 0, len(req.GetChildren()))
	for _, id := range req.GetChildren() {
		inheritors = append(inheritors, Permission{FileID: id, UserID: userID, Inherited: rootPermissionID})
	}

	_, err = s.CreateMany(ctx, inheritors)
	if err != nil {
		s.logger.Errorf("failed creating inheritors permissions: %v", err)
		return nil, err
	}

	return &pb.CreatePermissionResponse{Id: rootPermissionID.Hex()}, nil
}

// Create creates a permission of a file to a user,
// If successful returns the permission's ID and a nil error,
// otherwise returns empty string and non-nil error if any occured.
func (s Service) Create(ctx context.Context, permission Permission) (primitive.ObjectID, error) {
	collection := s.db.Collection(PermissionCollectionName)
	result, err := collection.InsertOne(ctx, permission)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// CreateMany creates many permissions using InsertMany.
func (s Service) CreateMany(ctx context.Context, permissions []Permission) ([]primitive.ObjectID, error) {
	collection := s.db.Collection(PermissionCollectionName)
	documents := make([]interface{}, 0, len(permissions))
	for _, permission := range permissions {
		documents = append(documents, permission)
	}

	results, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, 0, len(results.InsertedIDs))
	for _, result := range results.InsertedIDs {
		ids = append(ids, result.(primitive.ObjectID))
	}

	return ids, nil
}

// GetPermission is the request handler for retrieving permission details by its ID.
func (s Service) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.PermissionObject, error) {
	permissionID := req.GetId()
	permission, err := s.Get(ctx, permissionID)
	if err != nil {
		s.logger.Errorf(err.Error())
		return nil, err
	}

	var response *pb.PermissionObject

	if permission.ID != primitive.NilObjectID {
		response = &pb.PermissionObject{
			Id:     permission.ID.Hex(),
			FileID: permission.FileID,
			UserID: permission.UserID,
		}

		if permission.Inherited != primitive.NilObjectID {
			response.Inherited = permission.Inherited.Hex()
		}
	}

	return response, nil
}

// Get finds one permission that its ID matches permissionID,
// if successful returns the permission, and a nil error,
// if the permission is not found it would return a zero value Permission{},
// otherwise returns nil and non-nil error if any occured.
func (s Service) Get(ctx context.Context, permissionID string) (Permission, error) {
	collection := s.db.Collection(PermissionCollectionName)
	objectID, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return Permission{}, err
	}

	var permission Permission
	filter := bson.D{bson.E{Key: "_id", Value: objectID}}
	err = collection.FindOne(ctx, filter).Decode(&permission)
	if err != nil {
		return Permission{}, err
	}

	return permission, nil
}

// GetByFileAndUser finds one permission that its FileID and UserID matches fileID and userID,
// if successful returns the permission, and a nil error,
// if the permission is not found it would return a zero value Permission{},
// otherwise returns nil and non-nil error if any occured.
func (s Service) GetByFileAndUser(ctx context.Context, fileID string, userID string) (Permission, error) {
	collection := s.db.Collection(PermissionCollectionName)

	var permission Permission
	filter := bson.D{bson.E{Key: "fileID", Value: fileID}, bson.E{Key: "userID", Value: userID}}
	err := collection.FindOne(ctx, filter).Decode(&permission)
	if err != nil {
		return Permission{}, err
	}

	return permission, nil
}

// DeletePermission is the request handler for deleting permission by its ID.
func (s Service) DeletePermission(
	ctx context.Context, req *pb.DeletePermissionRequest,
) (*pb.PermissionObject, error) {
	permissionID := req.GetId()
	permission, err := s.Delete(ctx, permissionID)
	if err != nil {
		s.logger.Errorf(err.Error())
		return nil, err
	}

	var response *pb.PermissionObject

	if permission.ID != primitive.NilObjectID {
		response = &pb.PermissionObject{
			Id:     permission.ID.Hex(),
			FileID: permission.FileID,
			UserID: permission.UserID,
		}

		if permission.Inherited != primitive.NilObjectID {
			response.Inherited = permission.Inherited.Hex()
		}
	}

	return response, nil
}

// Delete finds one permission that its ID matches permissionID and deletes it,
// if successful returns the deleted permission, otherwise returns zero-value Permission,
// and non-nil error if any occured.
func (s Service) Delete(ctx context.Context, permissionID string) (Permission, error) {
	collection := s.db.Collection(PermissionCollectionName)
	inheritors, err := s.GetInheritors(ctx, permissionID)
	if err != nil {
		return Permission{}, err
	}

	for _, inheritor := range inheritors {
		filter := bson.D{bson.E{Key: "_id", Value: inheritor.ID}}
		result, err := collection.DeleteOne(ctx, filter)
		if err != nil || result.DeletedCount != 1 {
			return Permission{}, err
		}
	}

	objectID, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return Permission{}, err
	}

	var permission Permission
	filter := bson.D{bson.E{Key: "_id", Value: objectID}}
	err = collection.FindOneAndDelete(ctx, filter).Decode(&permission)
	if err != nil {
		return Permission{}, err
	}

	return permission, nil
}

// GetInheritors returns a slice of the inheritors of the permission that its ID matches permissionID.
func (s Service) GetInheritors(ctx context.Context, permissionID string) ([]Permission, error) {
	collection := s.db.Collection(PermissionCollectionName)

	inheritors := make([]Permission, 0)
	currentPermissionID := permissionID
	for {
		objectID, err := primitive.ObjectIDFromHex(currentPermissionID)
		if err != nil {
			return nil, err
		}

		// Find a permission which inherited from currentPermissionID.
		var childPermission Permission
		filter := bson.D{bson.E{Key: "inherited", Value: objectID}}
		err = collection.FindOne(ctx, filter).Decode(&childPermission)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				break
			}
		}

		inheritors = append(inheritors, childPermission)
		currentPermissionID = childPermission.ID.Hex()
	}

	return inheritors, nil
}
