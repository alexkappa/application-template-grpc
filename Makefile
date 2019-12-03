SHELL := /bin/bash

build: generate
	@go build -o bin/grpc-demo

generate:
	@go generate ./...

clean:
	@rm -rf \
		api/*/*.pb.go \
		api/*/*.pb.gw.go \
		api/*/*.swagger.json

.PHONY: build generate clean