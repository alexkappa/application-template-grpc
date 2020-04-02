SHELL := /bin/bash

build: generate
	@go build -o bin/service

generate:
	@go generate ./...

test:
	@go test ./...

clean:
	@rm -rf \
		api/*/*.pb.go \
		api/*/*.pb.gw.go \
		api/*/*.swagger.json

.PHONY: build generate clean