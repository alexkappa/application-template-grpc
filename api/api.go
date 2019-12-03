package api

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Service should be implemented by all gRPC services wishing to be exposed as
// resource oriented APIs.
type Service interface {

	// Register the service with gRPC gateway.
	Register(ctx context.Context, mux *runtime.ServeMux) error
}
