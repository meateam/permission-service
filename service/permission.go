package service

import (
	pb "github.com/meateam/permission-service/proto"
)

// Permission is an interface of a permission object.
type Permission interface {
	GetID() string

	GetFileID() string

	GetUserID() string

	MarshalProto(permission *pb.PermissionObject) error
}
