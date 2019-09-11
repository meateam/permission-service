package mongodb

import (
	"fmt"

	pb "github.com/meateam/permission-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BSON is the structure that represents a permission as it's stored.
type BSON struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FileID    string             `bson:"fileID,omitempty"`
	UserID    string             `bson:"userID,omitempty"`
	Inherited primitive.ObjectID `bson:"inherited,omitempty"`
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
	if b == nil {
		panic("b == nil")
	}

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
