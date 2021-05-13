package api

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Service should be implemented by all gRPC services wishing to be exposed as
// resource oriented APIs.
type Service interface {
	// Register the service with gRPC and gRPC-Gateway.
	Register(ctx context.Context, registrar grpc.ServiceRegistrar, handler *runtime.ServeMux) error
}
