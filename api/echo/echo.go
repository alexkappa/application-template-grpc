package echo

import (
	"context"

	"github.com/alexkappa/service-template-grpc/api"
	"github.com/alexkappa/service-template-grpc/api/echo/proto"
	"github.com/alexkappa/service-template-grpc/pkg/store"
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
func (s *echoService) Echo(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	v, err := s.store.Echo(req.Value)
	if err != nil {
		return nil, err
	}
	return &proto.EchoResponse{Value: v}, nil
}

// Register the service with gRPC gateway.
func (s *echoService) Register(ctx context.Context, mux *runtime.ServeMux) error {
	return proto.RegisterEchoHandlerServer(ctx, mux, s)
}

var _ proto.EchoServer = (*echoService)(nil)
var _ api.Service = (*echoService)(nil)
