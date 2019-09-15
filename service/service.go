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
func (s Service) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.PermissionObject, error) {
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

	permission, err := s.controller.CreatePermission(ctx, fileID, userID, role)
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
func (s Service) GetFilePermissions(ctx context.Context, req *pb.GetFilePermissionsRequest) (*pb.GetFilePermissionsResponse, error) {
	fileID := req.GetFileID()
	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	filePermissions, err := s.controller.GetFilePermissions(ctx, "")
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
		return nil, fmt.Errorf("UserID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("FileID is required")
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
	isPermitted, err := s.controller.IsPermitted(ctx, fileID, userID, role)
	if err != nil {
		return &pb.IsPermittedResponse{Permitted: false}, err
	}

	return &pb.IsPermittedResponse{Permitted: isPermitted}, nil
}
