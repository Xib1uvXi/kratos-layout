package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos-layout/api/health"
)

// HealthService is a health check service.
type HealthService struct {
	health.UnimplementedHealthServer
}

// NewHealthService creates a new health service.
func NewHealthService() *HealthService {
	return &HealthService{}
}

// Check implements health.HealthServer.
func (s *HealthService) Check(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
