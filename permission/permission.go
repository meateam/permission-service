package permission

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Service is a structure used for handling Permission Service grpc requests.
type Service struct {
	mongoClient *mongo.Client
	logger      *logrus.Logger
}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s *Service) HealthCheck(mongoClientPingTimeout time.Duration) bool {
	if s == nil {
		panic("s == nil")
	}

	timeoutCtx, cancel := context.WithTimeout(context.TODO(), mongoClientPingTimeout)
	defer cancel()
	err := s.mongoClient.Ping(timeoutCtx, readpref.Primary())
	if err != nil {
		s.logger.Errorf("%v", err)
		return false
	}

	return true
}

// NewService creates a Service and returns it.
func NewService(mongoClient *mongo.Client, logger *logrus.Logger) *Service {
	return &Service{mongoClient: mongoClient, logger: logger}
}
