# Application Template - gRPC

This is an application template for a Go service using gRPC and gRPC Gateway.

## Install

For the official instructions for installing `grpc` and `grpc-gateway` check out:

- https://grpc.io/docs/quickstart/go/
- https://github.com/grpc-ecosystem/grpc-gateway

## Conventions

Each service is defined in its own package under the `api` directory. For example to create a new Echo service, you would create a package `api/echo`.

This package should include a `echo.proto` detailing the protocol buffer definition of the service as well as any necessary gRPC gateway options to expose the service as an HTTP handler as well.

It should also contain an instruction to compile the protocol buffer definition to `*.pb.go`, `*.pb.gw.go` and `*.swagger.json` files. The convention is to add the following code into a `echo.go` file in the same package.

```Go
//go:generate bash ../../scripts/protoc.sh *.proto
package echo
```

With this in place, running `make generate` will generate all necessary files. 

Using the generated stubs, we must then implement the service and export it for use with  `api.Server`. So let's expand the `foo.go` file with the following.

```Go
//go:generate bash ../../scripts/protoc.sh *.proto
package echo

import (
	"context"

	"github.com/alexkappa/service-template-grpc/api"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type echoService struct{}

// Service exposes this service 
func Service() api.Service {
	return &echoService{}
}

// Echo handles an echo request.
func (s *echoService) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Value: req.Value}, nil
}

// Register the service with gRPC gateway. 
// 
// This function satisfies the `api.Service` interface.
func (s *echoService) Register(ctx context.Context, mux *runtime.ServeMux) error {
	return RegisterEchoHandlerServer(ctx, mux, s)
}

var _ EchoServer = (*echoService)(nil)
```

Notice how the `Register(...)` function satisfies the `api.Service` interface.