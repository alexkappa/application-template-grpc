package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	httpAddr    string
	httpHandler *runtime.ServeMux
	http        *http.Server
	grpcAddr    string
	grpc        *grpc.Server
	ctx         context.Context
	cancel      context.CancelFunc
	logger      *log.Logger
}

func NewServer(options ...ServerOption) *Server {
	s := &Server{
		httpAddr: ":11001",
		grpcAddr: ":11010",
		logger:   log.New(),
		ctx:      context.Background(),
	}
	s.logger.SetFormatter(&log.JSONFormatter{})
	for _, fn := range options {
		fn(s)
	}
	s.ctx, s.cancel = signal.NotifyContext(s.ctx, os.Interrupt)
	s.grpc = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_tags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(log.NewEntry(s.logger)),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(),
			grpc_tags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(log.NewEntry(s.logger)),
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

	err := make(chan error)

	go s.serveHttp(err)
	go s.serveRpc(err)

	select {
	case e := <-err:
		s.stopHttp()
		s.stopRpc()
		return e
	case <-s.ctx.Done():
		if errors.Is(s.ctx.Err(), context.Canceled) {
			s.logger.Print("SIGINT received")
			s.stopHttp()
			s.stopRpc()
			return nil
		}
		return s.ctx.Err()
	}
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
		s.logger.Fatalf("failed to listen: %v", err)
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
