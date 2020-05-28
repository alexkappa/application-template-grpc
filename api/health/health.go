package health

import (
	context "context"

	"github.com/alexkappa/service-template-grpc/api"
	"github.com/alexkappa/service-template-grpc/api/health/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type healthService struct{}

func Service() *healthService {
	return &healthService{}
}

func (s *healthService) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{Status: proto.HealthCheckResponse_SERVING}, nil
}

func (s *healthService) Register(ctx context.Context, mux *runtime.ServeMux) error {
	return proto.RegisterHealthHandlerServer(ctx, mux, s)
}

var _ proto.HealthServer = (*healthService)(nil)
var _ api.Service = (*healthService)(nil)
