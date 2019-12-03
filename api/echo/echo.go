//go:generate bash ../../scripts/protoc.sh *.proto
package echo

import (
	"context"

	"github.com/alexkappa/grpc-demo/api"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type echoService struct{}

func Service() *echoService { return new(echoService) }

// Echo handles an echo request.
func (s *echoService) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Value: req.Value}, nil
}

// Register the service with gRPC gateway.
func (s *echoService) Register(ctx context.Context, mux *runtime.ServeMux) error {
	return RegisterEchoHandlerServer(ctx, mux, s)
}

var _ EchoServer = (*echoService)(nil)
var _ api.Service = (*echoService)(nil)
