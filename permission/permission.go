package permission

import (
	"github.com/sirupsen/logrus"
)

// Service is a structure used for handling Permission Service grpc requests.
type Service struct{}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s *Service) HealthCheck() bool {
	return true
}

// NewService creates a Service and returns it.
func NewService(mongoClient interface{}, logger *logrus.Logger) *Service {
	return &Service{}
}
