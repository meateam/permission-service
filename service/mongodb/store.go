package mongodb

import (
	"context"
	"fmt"

	pb "github.com/meateam/permission-service/proto"
	"github.com/meateam/permission-service/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	// PermissionBSONRoleField is the name of the role field in BSON.
	PermissionBSONRoleField = "role"

	// PermissionBSONCreatorField is the name of the creator field in BSON.
	PermissionBSONCreatorField = "creator"

	// PermissionBSONAppIDField is the name of the appID field in BSON.
	PermissionBSONAppIDField = "appID"

	// PermissionBSONCreatedAtField is the name of the created_at field in BSON.
	PermissionBSONCreatedAtField = "created_at"

	// PermissionBSONUpdatedAtField is the name of the updated_at field in BSON.
	PermissionBSONUpdatedAtField = "updated_at"
)

// MongoStore holds the mongodb database and implements Store interface.
type MongoStore struct {
	DB *mongo.Database
}

// PagingRes is the struct of the result from the pagination requests
type PagingRes struct {
	permissions []service.Permission
	itemCount   int64
	pageNum     int64
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
// If permission already exists then it's updated to have permission values,
// If successful returns the permission and a nil error,
// Override indicates whether to update the permission if already exists, or not and return error.
// otherwise returns empty string and non-nil error if any occurred.
func (s MongoStore) Create(
	ctx context.Context,
	permission service.Permission,
	override bool,
) (service.Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)
	fileID := permission.GetFileID()
	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	userID := permission.GetUserID()
	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	role := permission.GetRole()
	if pb.Role_name[int32(role)] == "" {
		return nil, fmt.Errorf("role does not exist")
	}

	creator := permission.GetCreator()
	if userID == "" {
		return nil, fmt.Errorf("creator is required")
	}

	appID := permission.GetAppID()
	if appID == "" {
		return nil, fmt.Errorf("appID is required")
	}

	createdAt := permission.GetCreatedAt()
	if appID == "" {
		return nil, fmt.Errorf("current time could not be obtained")
	}

	updatedAt := permission.GetUpdatedAt()
	if appID == "" {
		return nil, fmt.Errorf("current time could not be obtained")
	}

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

	// Check existance of the permission. If the error received is not `404` then abort.
	// In case the permission exists, do not update createdAt.
	existingPermission, existingPermissionErr := s.Get(ctx, filter)
	if existingPermissionErr != nil && existingPermissionErr != mongo.ErrNoDocuments {
		return nil, existingPermissionErr
	} else if existingPermissionErr == nil {
		createdAt = existingPermission.GetCreatedAt()
	}

	newPermission := bson.D{
		bson.E{
			Key:   PermissionBSONFileIDField,
			Value: fileID,
		},
		bson.E{
			Key:   PermissionBSONUserIDField,
			Value: userID,
		},
		bson.E{
			Key:   PermissionBSONRoleField,
			Value: role,
		},
		bson.E{
			Key:   PermissionBSONCreatorField,
			Value: creator,
		},
		bson.E{
			Key:   PermissionBSONAppIDField,
			Value: appID,
		},
		bson.E{
			Key:   PermissionBSONCreatedAtField,
			Value: createdAt,
		},
		bson.E{
			Key:   PermissionBSONUpdatedAtField,
			Value: updatedAt,
		},
	}

	update := bson.D{
		bson.E{
			Key:   "$set",
			Value: newPermission,
		},
	}

	// In case override is false, check if there is a permission, and if there is one - return it.
	if !override && existingPermissionErr == nil {
		return existingPermission, nil
	}

	// If override is true, or false and there is no permission existing,
	// then update and allow to override the permission fields
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, opts)

	updatedPermission := &BSON{}
	err := result.Decode(updatedPermission)
	if err != nil {
		return nil, err
	}

	return updatedPermission, nil
}

// Get finds one permission that matches filter,
// if successful returns the permission, and a nil error,
// if the permission is not found it would return nil and NotFound error,
// otherwise returns nil and non-nil error if any occurred.
func (s MongoStore) Get(ctx context.Context, filter interface{}) (service.Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)

	permission := &BSON{}
	foundPerm := collection.FindOne(ctx, filter)
	err := foundPerm.Decode(permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// GetAll finds all permissions that matches filter,
// if successful returns the permissions, and a nil error,
// otherwise returns nil and non-nil error if any occurred.
func (s MongoStore) GetAll(ctx context.Context, filter interface{}) ([]service.Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	permissions := []service.Permission{}
	for cur.Next(ctx) {
		permission := &BSON{}
		err := cur.Decode(permission)
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

// GetUserPermissionsByPage returns a slice of the permissions requested by the filter,
// in the page they belong to by page number and page size.
// sortBy is the field by which the sorting of the permissions will be committed,
// defaulting to reverse mongoID.
func (s MongoStore) GetUserPermissionsByPage(ctx context.Context, pn int64, ps int64, sortBy bson.D, filter interface{}) (*PagingRes, error) {
	collection := s.DB.Collection(PermissionCollectionName)

	itemCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(sortBy).SetLimit(ps).SetSkip(pn * ps)

	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	permissions := []service.Permission{}
	for cur.Next(ctx) {
		permission := &BSON{}
		err := cur.Decode(permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	var pr *PagingRes = &PagingRes{permissions: permissions, pageNum: pn, itemCount: itemCount}

	return pr, nil
}

// Delete finds the first permission that matches filter and deletes it,
// if successful returns the deleted permission, otherwise returns nil,
// and non-nil error if any occurred.
func (s MongoStore) Delete(ctx context.Context, filter interface{}) (service.Permission, error) {
	collection := s.DB.Collection(PermissionCollectionName)
	permission := &BSON{}
	if err := collection.FindOneAndDelete(ctx, filter).Decode(permission); err != nil {
		return nil, err
	}

	return permission, nil
}
