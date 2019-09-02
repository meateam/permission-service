package permission

import (
	"context"
	"time"

	pb "github.com/meateam/permission-service/proto"
	"github.com/sirupsen/logrus"
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
	FileID    *string            `json:"fileID,omitempty" bson:"fileID,omitempty"`
	UserID    *string            `json:"userID,omitempty" bson:"userID,omitempty"`
	Inherited *string            `json:"inherited,omitempty" bson:"inherited,omitempty"`
}

// Store is an interface for handling the storing of permissions.
type Store interface {
	Create(ctx context.Context, permission Permission) (string, error)
	CreateMany(ctx context.Context, permissions []Permission) ([]string, error)
	Get(ctx context.Context, permissionID string) (Permission, error)
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

	rootPermission := Permission{FileID: &rootFileID, UserID: &userID}
	// Create the root permission.
	rootPermissionID, err := s.Create(ctx, rootPermission)
	if err != nil {
		s.logger.Error("failed creating root permission: %v", err)
		return nil, err
	}

	inheritors := make([]Permission, 0, len(req.GetChildren()))
	for _, id := range req.GetChildren() {
		inheritors = append(inheritors, Permission{FileID: &id, UserID: &userID, Inherited: &rootFileID})
	}

	_, err = s.CreateMany(ctx, inheritors)
	if err != nil {
		s.logger.Errorf("failed creating inheritors permissions: %v", err)
		return nil, err
	}

	return &pb.CreatePermissionResponse{Id: rootPermissionID}, nil
}

// Create creates a permission of a file to a user,
// If successful returns the permission's ID and a nil error,
// otherwise returns empty string and non-nil error if any occured.
func (s Service) Create(ctx context.Context, permission Permission) (string, error) {
	collection := s.db.Collection(PermissionCollectionName)
	result, err := collection.InsertOne(ctx, permission)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// CreateMany creates many permissions using InsertMany.
func (s Service) CreateMany(ctx context.Context, permissions []Permission) ([]string, error) {
	collection := s.db.Collection(PermissionCollectionName)
	documents := make([]interface{}, 0, len(permissions))
	for _, permission := range permissions {
		documents = append(documents, permission)
	}

	results, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(results.InsertedIDs))
	for _, result := range results.InsertedIDs {
		ids = append(ids, result.(primitive.ObjectID).Hex())
	}

	return ids, nil
}
