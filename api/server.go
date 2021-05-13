package api

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	akctx "github.com/alexkappa/service-template-grpc/pkg/context"
)

type Server struct {
	httpAddr    string
	httpHandler *runtime.ServeMux
	http        *http.Server
	grpcAddr    string
	grpc        *grpc.Server
	ctx         context.Context
	logger      *log.Logger
}

func NewServer(options ...ServerOption) *Server {
	s := &Server{
		httpAddr: ":11001",
		grpcAddr: ":11010",
		ctx:      context.Background(),
		logger:   log.New(),
	}
	s.logger.SetFormatter(&log.JSONFormatter{})
	for _, fn := range options {
		fn(s)
	}
	s.grpc = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			unaryInterceptorRequestId(),
			unaryInterceptorlogging(s.logger),
		),
	)
	s.httpHandler = runtime.NewServeMux()
	s.http = &http.Server{
		Addr:    s.httpAddr,
		Handler: s.httpHandler,
	}
	return s
}

func (s *Server) Register(services ...Service) {
	for _, service := range services {
		service.Register(s.ctx, s.grpc, s.httpHandler)
	}
}

func (s *Server) Serve() error {
	hupCh := make(chan os.Signal, 1)
	signal.Notify(hupCh, os.Interrupt)

	errCh := make(chan error)

	go s.serveHttp(errCh)
	go s.serveRpc(errCh)

	select {
	case err := <-errCh:
		s.stopHttp()
		s.stopRpc()
		return err
	case <-hupCh:
		s.logger.Print("HUP signal received")
		s.stopHttp()
		s.stopRpc()
	}

	return nil
}

func (s *Server) serveHttp(errCh chan<- error) {
	s.logger.Printf("HTTP server listening on %s", s.httpAddr)
	err := s.http.ListenAndServe()
	if err != nil {
		errCh <- err
	}
}

func (s *Server) stopHttp() {
	s.http.Close()
	s.logger.Print("HTTP server stopped")
}

func (s *Server) serveRpc(errCh chan<- error) {
	l, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.logger.Printf("gRPC server listening on %s", s.grpcAddr)
	err = s.grpc.Serve(l)
	if err != nil {
		errCh <- err
	}
}

func (s *Server) stopRpc() {
	s.grpc.Stop()
	s.logger.Print("gRPC server stopped")
}

type ServerOption func(*Server)

func WithHTTPAddress(addr string) ServerOption {
	return func(s *Server) {
		s.httpAddr = addr
	}
}

func WithRPCAddress(addr string) ServerOption {
	return func(s *Server) {
		s.grpcAddr = addr
	}
}

func WithContext(ctx context.Context) ServerOption {
	return func(s *Server) {
		s.ctx = ctx
	}
}

func WithLogger(logger *log.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

func unaryInterceptorRequestId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(akctx.WithRequestID(ctx), req)
	}
}

func unaryInterceptorlogging(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		res, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		logger.
			WithField("req_id", akctx.RequestID(ctx)).
			WithField("req", req).
			WithField("res", res).
			WithField("method", info.FullMethod).
			Print("remote procedure call")

		return res, err
	}
}
