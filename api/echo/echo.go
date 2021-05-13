package echo

import (
	"context"
	"errors"

	"github.com/alexkappa/service-template-grpc/api"
	"github.com/alexkappa/service-template-grpc/pkg/store"
	proto "github.com/alexkappa/service-template-grpc/proto/echo/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type echoService struct {
	store store.KVStore
}

// Service returns a new instance of the echo service.
//
// For illustration purposes a `store.EchoStore` is passed as an argument to the
// Service.
func Service(s store.KVStore) api.Service {
	return &echoService{s}
}

// Echo handles an echo request.
func (s *echoService) Echo(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	c, err := s.store.Get(ctx, req.Value)
	if err != nil && errors.Is(err, store.ErrKeyNotFount) {
		c = int64(0)
	}
	c = c.(int64) + 1
	s.store.Set(ctx, req.Value, c)
	return &proto.EchoResponse{Value: req.Value, Count: c.(int64)}, nil
}

// Register the service.
func (s *echoService) Register(ctx context.Context, registrar grpc.ServiceRegistrar, handler *runtime.ServeMux) error {
	proto.RegisterEchoServer(registrar, s)
	proto.RegisterEchoHandlerServer(ctx, handler, s)
	return nil
}

func Client(c grpc.ClientConnInterface) proto.EchoClient {
	return proto.NewEchoClient(c)
}

var _ proto.EchoServer = (*echoService)(nil)
var _ api.Service = (*echoService)(nil)
