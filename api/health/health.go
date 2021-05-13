package health

import (
	"context"

	"github.com/alexkappa/service-template-grpc/api"
	proto "github.com/alexkappa/service-template-grpc/proto/health/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type healthService struct{}

// Service returns a new instance of the health service.
func Service() api.Service {
	return &healthService{}
}

func (s *healthService) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{Status: proto.HealthCheckResponse_SERVING}, nil
}

// Register the service with gRPC gateway.
func (s *healthService) Register(ctx context.Context, registrar grpc.ServiceRegistrar, handler *runtime.ServeMux) error {
	proto.RegisterHealthServer(registrar, s)
	proto.RegisterHealthHandlerServer(ctx, handler, s)
	return nil
}

func Client(c grpc.ClientConnInterface) proto.HealthClient {
	return proto.NewHealthClient(c)
}

var _ proto.HealthServer = (*healthService)(nil)
var _ api.Service = (*healthService)(nil)
