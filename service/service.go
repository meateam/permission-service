package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/meateam/permission-service/proto"
	"github.com/sirupsen/logrus"
)

// Service is a structure used for handling Permission Service grpc requests.
type Service struct {
	controller Controller
	logger     *logrus.Logger
}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s Service) HealthCheck(mongoClientPingTimeout time.Duration) bool {
	timeoutCtx, cancel := context.WithTimeout(context.TODO(), mongoClientPingTimeout)
	defer cancel()
	healthy, err := s.controller.HealthCheck(timeoutCtx)
	if err != nil {
		s.logger.Errorf("%v", err)
		return false
	}

	return healthy
}

// NewService creates a Service and returns it.
func NewService(controller Controller, logger *logrus.Logger) Service {
	return Service{controller: controller, logger: logger}
}

// CreatePermission is the request handler for creating a permission of a file to user.
func (s Service) CreatePermission(
	ctx context.Context,
	req *pb.CreatePermissionRequest,
) (*pb.PermissionObject, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()
	role := req.GetRole()
	creator := req.GetCreator()
	override := req.GetOverride()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	if pb.Role_name[int32(role)] == "" {
		return nil, fmt.Errorf("role does not exist")
	}

	if creator == "" {
		return nil, fmt.Errorf("creator is required")
	}

	permission, err := s.controller.CreatePermission(ctx, fileID, userID, role, creator, override)
	if err != nil {
		return nil, err
	}

	var response pb.PermissionObject
	if err := permission.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetFilePermissions is the request handler for retrieving permissions of file by its ID.
func (s Service) GetFilePermissions(
	ctx context.Context,
	req *pb.GetFilePermissionsRequest,
) (*pb.GetFilePermissionsResponse, error) {
	fileID := req.GetFileID()
	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	filePermissions, err := s.controller.GetFilePermissions(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return &pb.GetFilePermissionsResponse{Permissions: filePermissions}, nil
}

// DeletePermission is the request handler for deleting permission by its ID.
func (s Service) DeletePermission(
	ctx context.Context, req *pb.DeletePermissionRequest,
) (*pb.PermissionObject, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	permission, err := s.controller.DeletePermission(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	var response pb.PermissionObject
	if err = permission.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetPermission is the request handler for retrieving a permission by a user and file ids.
func (s Service) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.PermissionObject, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()
	if userID == "" {
		return nil, fmt.Errorf("UserID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("FileID is required")
	}

	permission, err := s.controller.GetByFileAndUser(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	var response pb.PermissionObject
	if err = permission.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// IsPermitted is the request handler for checking user permission by userID and fileID.
func (s Service) IsPermitted(ctx context.Context, req *pb.IsPermittedRequest) (*pb.IsPermittedResponse, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()
	role := req.GetRole()
	if userID == "" {
		return nil, fmt.Errorf("UserID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("FileID is required")
	}

	if pb.Role_name[int32(role)] == "" {
		return nil, fmt.Errorf("role does not exist")
	}

	permission, err := s.controller.GetByFileAndUser(ctx, fileID, userID)
	if err != nil {
		return &pb.IsPermittedResponse{Permitted: false}, err
	}

	isPermitted := isSubRole(permission.GetRole(), role)
	return &pb.IsPermittedResponse{Permitted: isPermitted}, nil
}

// GetUserPermissions is the request handler for fetching the permissions that a user has.
func (s Service) GetUserPermissions(
	ctx context.Context,
	req *pb.GetUserPermissionsRequest) (*pb.GetUserPermissionsResponse, error) {
	userID := req.GetUserID()
	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	permissions, err := s.controller.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserPermissionsResponse{Permissions: permissions}, nil
}

// DeleteFilePermissions is the request handler for deleting all permissions that exist for a certain file.
func (s Service) DeleteFilePermissions(
	ctx context.Context,
	req *pb.DeleteFilePermissionsRequest,
) (*pb.DeleteFilePermissionsResponse, error) {
	fileID := req.GetFileID()
	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	permissions, err := s.controller.DeleteFilePermissions(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteFilePermissionsResponse{Permissions: permissions}, nil
}

func isSubRole(role pb.Role, wanted pb.Role) bool {
	if wanted == pb.Role_NONE {
		return false
	}

	switch role {
	case pb.Role_NONE:
		return false
	case pb.Role_WRITE:
		return wanted == pb.Role_WRITE || wanted == pb.Role_READ
	case pb.Role_READ:
		return wanted == pb.Role_READ
	default:
		return false
	}
}
