package api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	addr  string
	ctx   context.Context
	mux   *runtime.ServeMux
	mongo *mongo.Client
}

func NewServer(options ...ServerOption) *Server {
	s := &Server{}
	s.mux = runtime.NewServeMux()
	for _, fn := range options {
		fn(s)
	}
	return s
}

func (s *Server) Register(services ...Service) {
	for _, service := range services {
		service.Register(s.ctx, s.mux)
	}
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.addr, s.mux)
}

type ServerOption func(*Server)

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithContext(ctx context.Context) ServerOption {
	return func(s *Server) {
		s.ctx = ctx
	}
}

func WithMongo(m *mongo.Client) ServerOption {
	return func(s *Server) {
		s.mongo = m
	}
}
