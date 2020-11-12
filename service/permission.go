package service

import (
	"time"

	pb "github.com/meateam/permission-service/proto"
)

// Permission is an interface of a permission object.
type Permission interface {
	GetID() string

	SetID(id string) error

	GetFileID() string

	SetFileID(fileID string) error

	GetUserID() string

	SetUserID(userID string) error

	GetRole() pb.Role

	SetRole(role pb.Role) error

	GetCreator() string

	SetCreator(creator string) error

	GetCreatedAt() string

	SetCreatedAt(createdAt time.Time) error

	GetUpdatedAt() string

	SetUpdatedAt(updatedAt time.Time) error

	MarshalProto(permission *pb.PermissionObject) error
}
