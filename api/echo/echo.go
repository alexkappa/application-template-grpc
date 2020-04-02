//go:generate bash ../../scripts/protoc.sh *.proto
package echo

import (
	"context"

	"github.com/alexkappa/application-template-grpc/api"
	"github.com/alexkappa/application-template-grpc/pkg/store"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type echoService struct {
	store store.EchoStore
}

// Service returns a new instance of the echo service.
//
// For illustration purposes a `store.EchoStore` is passed as an argument to the
// Service.
func Service(s store.EchoStore) api.Service {
	return &echoService{s}
}

// Echo handles an echo request.
func (s *echoService) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	v, err := s.store.Echo(req.Value)
	if err != nil {
		return nil, err
	}
	return &EchoResponse{Value: v}, nil
}

// Register the service with gRPC gateway.
func (s *echoService) Register(ctx context.Context, mux *runtime.ServeMux) error {
	return RegisterEchoHandlerServer(ctx, mux, s)
}

var _ EchoServer = (*echoService)(nil)
var _ api.Service = (*echoService)(nil)
