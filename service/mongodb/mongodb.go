package mongodb

import (
	"fmt"
	"time"

	pb "github.com/meateam/permission-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BSON is the structure that represents a permission as it's stored.
type BSON struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FileID    string             `bson:"fileID,omitempty"`
	UserID    string             `bson:"userID,omitempty"`
	Role      pb.Role            `bson:"role"`
	Creator   string             `bson:"creator"`
	AppID     string             `bson:"appID"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt,omitempty"`
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

// GetRole returns b.Role.
func (b BSON) GetRole() pb.Role {
	return b.Role
}

// SetRole sets b.Role to role.
func (b *BSON) SetRole(role pb.Role) error {
	if b == nil {
		panic("b == nil")
	}

	if pb.Role_name[int32(role)] == "" {
		return fmt.Errorf("Role does not exist")
	}

	b.Role = role
	return nil
}

// GetCreator returns b.Creator.
func (b BSON) GetCreator() string {
	return b.Creator
}

// SetCreator sets b.Creator to creator.
func (b *BSON) SetCreator(creator string) error {
	if b == nil {
		panic("b == nil")
	}

	if creator == "" {
		return fmt.Errorf("Creator is required")
	}

	b.Creator = creator
	return nil
}

// GetAppID returns b.AppID.
func (b BSON) GetAppID() string {
	return b.AppID
}

// SetAppID sets b.AppID to appID.
func (b *BSON) SetAppID(appID string) error {
	if b == nil {
		panic("b == nil")
	}

	if appID == "" {
		return fmt.Errorf("appID is required")
	}

	b.AppID = appID
	return nil
}

// GetCreatedAt returns b.CreatedAt.
func (b BSON) GetCreatedAt() time.Time {
	return b.CreatedAt
}

// GetUpdatedAt returns b.UpdatedAt.
func (b BSON) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}

// MarshalProto marshals b into a permission.
func (b BSON) MarshalProto(permission *pb.PermissionObject) error {
	permission.Id = b.GetID()
	permission.FileID = b.GetFileID()
	permission.UserID = b.GetUserID()
	permission.Role = b.GetRole()
	permission.Creator = b.GetCreator()
	permission.AppID = b.GetAppID()
	permission.CreatedAt = b.GetCreatedAt().Unix()
	permission.UpdatedAt = b.GetUpdatedAt().Unix()

	return nil
}
