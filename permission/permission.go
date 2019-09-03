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
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// MongoObjectIDField is the default mongodb unique key.
	MongoObjectIDField = "_id"

	// PermissionCollectionName is the name of the permissions collection.
	PermissionCollectionName = "permissions"

	// PermissionBSONFileIDField is the name of the fileID field in BSON.
	PermissionBSONFileIDField = "fileID"

	// PermissionBSONUserIDField is the name of the userID field in BSON.
	PermissionBSONUserIDField = "userID"

	// PermissionBSONInheritedField is the name of the inherited field in BSON.
	PermissionBSONInheritedField = "inherited"
)

// Store is an interface for handling the storing of permissions.
type Store interface {
	Create(ctx context.Context, permission Permission) (string, error)
	CreateMany(ctx context.Context, permissions []Permission) ([]string, error)
	Get(ctx context.Context, filter interface{}) (Permission, error)
	GetAll(ctx context.Context, filter interface{}) ([]Permission, error)
	Delete(ctx context.Context, filter interface{}) (Permission, error)
	HealthCheck(ctx context.Context) (bool, error)
}

// Service is a structure used for handling Permission Service grpc requests.
type Service struct {
	store  Store
	logger *logrus.Logger
}

// MongoStore holds the mongodb database and implements Store interface.
type MongoStore struct {
	DB *mongo.Database
}

// Permission is an interface of a permission object.
type Permission interface {
	GetID() string

	GetFileID() string

	GetUserID() string

	GetInherited() string

	MarshalProto(permission *pb.PermissionObject) error
}

// BSON is the structure that represents a permission as it's stored.
type BSON struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FileID    string             `json:"fileID,omitempty" bson:"fileID,omitempty"`
	UserID    string             `json:"userID,omitempty" bson:"userID,omitempty"`
	Inherited primitive.ObjectID `json:"inherited,omitempty" bson:"inherited,omitempty"`
}

// GetID returns the string value of the b.ID.
func (b BSON) GetID() string {
	if b.ID.IsZero() {
		return ""
	}

	return b.ID.Hex()
}

// SetID sets the b.ID ObjectID's string value to id.
func (b *BSON) SetID(id string) error {
	if b == nil {
		panic("b == nil")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	b.ID = objectID
	return nil
}

// GetFileID returns b.FileID.
func (b BSON) GetFileID() string {
	return b.FileID
}

// SetFileID sets b.FileID to fileID.
func (b *BSON) SetFileID(fileID string) error {
	if b == nil {
		panic("b == nil")
	}

	if fileID == "" {
		return fmt.Errorf("FileID is required")
	}

	b.FileID = fileID
	return nil
}

// GetUserID returns b.UserID.
func (b BSON) GetUserID() string {
	return b.UserID
}

// SetUserID sets b.UserID to userID.
func (b *BSON) SetUserID(userID string) error {
	if userID == "" {
		return fmt.Errorf("UserID is required")
	}

	b.UserID = userID
	return nil
}

// GetInherited returns the string value of b.Inherited.
func (b BSON) GetInherited() string {
	if b.Inherited.IsZero() {
		return ""
	}

	return b.Inherited.Hex()
}

// SetInherited sets the b.Inherited ObjectID's string value to inherited.
func (b *BSON) SetInherited(inherited string) error {
	if b == nil {
		panic("b == nil")
	}

	objectID, err := primitive.ObjectIDFromHex(inherited)
	if err != nil {
		return err
	}

	b.Inherited = objectID
	return nil
}

// MarshalProto marshals b into a permission.
func (b BSON) MarshalProto(permission *pb.PermissionObject) error {
	permission.Id = b.GetID()
	permission.FileID = b.GetFileID()
	permission.UserID = b.GetUserID()

	if b.GetInherited() != "" {
		permission.Inherited = b.GetInherited()
	}

	return nil
}

// NewMongoStore returns a new MongoStore.
func NewMongoStore(db *mongo.Database) (MongoStore, error) {
	collection := db.Collection(PermissionCollectionName)
	indexes := collection.Indexes()
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			bson.E{
				Key:   PermissionBSONFileIDField,
				Value: 1,
			},
			bson.E{
				Key:   PermissionBSONUserIDField,
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := indexes.CreateOne(context.Background(), indexModel)
	if err != nil {
		return MongoStore{}, err
	}

	return MongoStore{DB: db}, nil
}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s MongoStore) HealthCheck(ctx context.Context) (bool, error) {
	err := s.DB.Client().Ping(ctx, readpref.Primary())
	if err != nil {
		return false, err
	}

	return true, nil
}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s *Service) HealthCheck(mongoClientPingTimeout time.Duration) bool {
	if s == nil {
		panic("s == nil")
	}

	timeoutCtx, cancel := context.WithTimeout(context.TODO(), mongoClientPingTimeout)
	defer cancel()
	healthy, err := s.store.HealthCheck(timeoutCtx)
	if err != nil {
		s.logger.Errorf("%v", err)
		return false
	}

	return healthy
}

// NewService creates a Service and returns it.
func NewService(store Store, logger *logrus.Logger) *Service {
	return &Service{store: store, logger: logger}
}

// CreatePermission is the request handler for creating a permission of a file to user.
func (s Service) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.CreatePermissionResponse, error) {
	rootFileID := req.GetFileID()
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("UserID is required")
	}

	if rootFileID == "" {
		return nil, fmt.Errorf("FileID is required")
	}

	// Create the root permission.
	rootPermission := BSON{FileID: rootFileID, UserID: userID}
	rootPermissionID, err := s.store.Create(ctx, rootPermission)
	if err != nil {
		return nil, fmt.Errorf("failed creating root permission: %v", err)
	}

	// Map inheritors to slice of Permission.
	inheritors := make([]Permission, 0, len(req.GetInheritors()))
	for _, id := range req.GetInheritors() {
		inheritor := &BSON{}
		if err := inheritor.SetFileID(id); err != nil {
			return nil, err
		}

		if err := inheritor.SetUserID(userID); err != nil {
			return nil, err
		}

		if err := inheritor.SetInherited(rootPermissionID); err != nil {
			return nil, err
		}

		inheritors = append(inheritors, *inheritor)
	}

	// Create the inheritors' permissions.
	_, err = s.store.CreateMany(ctx, inheritors)
	if err != nil {
		return nil, fmt.Errorf("failed creating inheritors permissions: %v", err)
	}

	return &pb.CreatePermissionResponse{Id: rootPermissionID}, nil
}

// Create creates a permission of a file to a user,
// If successful returns the permission's ID and a nil error,
// otherwise returns empty string and non-nil error if any occured.
func (s MongoStore) Create(ctx context.Context, permission Permission) (string, error) {
	collection := s.DB.Collection(PermissionCollectionName)
	result, err := collection.InsertOne(ctx, permission)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// CreateMany creates many permissions using InsertMany.
func (s MongoStore) CreateMany(ctx context.Context, permissions []Permission) ([]string, error) {
	collection := s.DB.Collection(PermissionCollectionName)
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

// GetPermission is the request handler for retrieving permission details by its ID.
func (s Service) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.PermissionObject, error) {
	permissionID := req.GetId()
	if permissionID == "" {
		return nil, fmt.Errorf("ID is required")
	}

	permissionObjectID, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		bson.E{
			Key:   MongoObjectIDField,
			Value: permissionObjectID,
		},
	}

	permission, err := s.store.Get(ctx, filter)
	if err != nil {
		return nil, err
	}

	var response pb.PermissionObject
	if err = permission.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Get finds one permission that matches filter,
// if successful returns the permission, and a nil error,
// if the permission is not found it would return a zero value BSON{},
// otherwise returns nil and non-nil error if any occured.
func (s MongoStore) Get(ctx context.Context, filter interface{}) (Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)

	permission := BSON{}
	err := collection.FindOne(ctx, filter).Decode(&permission)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.Unimplemented, "permission not found")
	}

	return permission, nil
}

// GetAll finds all permissions that matches filter,
// if successful returns the permission, and a nil error,
// if the permission is not found it would return a zero value BSON{},
// otherwise returns nil and non-nil error if any occured.
func (s MongoStore) GetAll(ctx context.Context, filter interface{}) ([]Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)

	cur, err := collection.Find(ctx, filter)
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	permissions := make([]Permission, 0, 1)
	for cur.Next(ctx) {
		permission := BSON{}
		err := cur.Decode(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// DeletePermission is the request handler for deleting permission by its ID.
func (s Service) DeletePermission(
	ctx context.Context, req *pb.DeletePermissionRequest,
) (*pb.PermissionObject, error) {
	permissionID := req.GetId()
	if permissionID == "" {
		return nil, fmt.Errorf("ID is required")
	}

	permissionObjectID, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		bson.E{
			Key:   MongoObjectIDField,
			Value: permissionObjectID,
		},
	}

	permission, err := s.store.Delete(ctx, filter)
	if err != nil {
		return nil, err
	}

	var response pb.PermissionObject
	if err = permission.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete finds the first permission that matches filter and deletes it,
// if successful returns the deleted permission, otherwise returns nil,
// and non-nil error if any occured.
func (s MongoStore) Delete(ctx context.Context, filter interface{}) (Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)
	inheritors, err := s.GetInheritors(ctx, filter)
	if err != nil {
		return nil, err
	}

	for _, inheritor := range inheritors {
		inheritorID, err := primitive.ObjectIDFromHex(inheritor.GetID())
		if err != nil {
			return nil, err
		}

		filter := bson.D{
			bson.E{
				Key:   "_id",
				Value: inheritorID,
			},
		}
		result, err := collection.DeleteOne(ctx, filter)
		if err != nil || result.DeletedCount != 1 {
			return nil, err
		}
	}

	permission := BSON{}
	if err := collection.FindOneAndDelete(ctx, filter).Decode(&permission); err != nil {
		return nil, err
	}

	return permission, nil
}

// GetInheritors returns a slice of the inheritors of the permission that matches filter.
func (s MongoStore) GetInheritors(ctx context.Context, filter interface{}) ([]Permission, error) {
	permission, err := s.Get(ctx, filter)
	if err != nil {
		return nil, err
	}

	permissionID, err := primitive.ObjectIDFromHex(permission.GetID())
	if err != nil {
		return nil, err
	}

	// Find permissions which inherited from permission.
	inheritedFilter := bson.D{
		bson.E{
			Key:   PermissionBSONInheritedField,
			Value: permissionID,
		},
	}
	inheritors, err := s.GetAll(ctx, inheritedFilter)
	if err != nil {
		return nil, err
	}

	return inheritors, nil
}
