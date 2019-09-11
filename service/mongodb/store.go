package mongodb

import (
	"context"

	"github.com/meateam/permission-service/service"
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
)

// MongoStore holds the mongodb database and implements Store interface.
type MongoStore struct {
	DB *mongo.Database
}

// newMongoStore returns a new store.
func newMongoStore(db *mongo.Database) (MongoStore, error) {
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
	if err := s.DB.Client().Ping(ctx, readpref.Primary()); err != nil {
		return false, err
	}

	return true, nil
}

// Create creates a permission of a file to a user,
// If successful returns the permission's ID and a nil error,
// otherwise returns empty string and non-nil error if any occured.
func (s MongoStore) Create(ctx context.Context, permission service.Permission) (string, error) {
	collection := s.DB.Collection(PermissionCollectionName)
	result, err := collection.InsertOne(ctx, permission)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// CreateMany creates many permissions using InsertMany.
func (s MongoStore) CreateMany(ctx context.Context, permissions []service.Permission) ([]string, error) {
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

// Get finds one permission that matches filter,
// if successful returns the permission, and a nil error,
// if the permission is not found it would return a zero value BSON{},
// otherwise returns nil and non-nil error if any occured.
func (s MongoStore) Get(ctx context.Context, filter interface{}) (service.Permission, error) {
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
func (s MongoStore) GetAll(ctx context.Context, filter interface{}) ([]service.Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)

	cur, err := collection.Find(ctx, filter)
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	permissions := make([]service.Permission, 0, 1)
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

// Delete finds the first permission that matches filter and deletes it,
// if successful returns the deleted permission, otherwise returns nil,
// and non-nil error if any occured.
func (s MongoStore) Delete(ctx context.Context, filter interface{}) (service.Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)
	permission := BSON{}
	if err := collection.FindOneAndDelete(ctx, filter).Decode(&permission); err != nil {
		return nil, err
	}

	return permission, nil
}
