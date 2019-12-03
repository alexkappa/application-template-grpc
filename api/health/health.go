//go:generate bash ../../scripts/protoc.sh *.proto
package health

import (
	context "context"

	"github.com/alexkappa/grpc-demo/api"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type healthService struct{}

func Service() *healthService {
	return &healthService{}
}

func (s *healthService) Check(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{Status: HealthCheckResponse_SERVING}, nil
}

func (s *healthService) Register(ctx context.Context, mux *runtime.ServeMux) error {
	return RegisterHealthHandlerServer(ctx, mux, s)
}

var _ HealthServer = (*healthService)(nil)
var _ api.Service = (*healthService)(nil)
