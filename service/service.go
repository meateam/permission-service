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
func (s Service) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.CreatePermissionResponse, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("UserID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("FileID is required")
	}

	id, err := s.controller.CreatePermission(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePermissionResponse{Id: id}, nil
}

// GetPermission is the request handler for retrieving permission details by its ID.
func (s Service) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.PermissionObject, error) {
	permissionID := req.GetId()
	if permissionID == "" {
		return nil, fmt.Errorf("ID is required")
	}

	permission, err := s.controller.GetPermission(ctx, "", "")
	if err != nil {
		return nil, err
	}

	var response pb.PermissionObject
	if err := permission.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeletePermission is the request handler for deleting permission by its ID.
func (s Service) DeletePermission(
	ctx context.Context, req *pb.DeletePermissionRequest,
) (*pb.PermissionObject, error) {
	permissionID := req.GetId()
	if permissionID == "" {
		return nil, fmt.Errorf("ID is required")
	}

	permission, err := s.controller.DeletePermission(ctx, "", "")
	if err != nil {
		return nil, err
	}

	var response pb.PermissionObject
	if err = permission.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
