SHELL := /bin/bash
PROTO_FILES := $(shell ls api/proto/*.proto api/*/proto/*.proto)

build: generate
	@go build -o bin/service

generate:
	@buf generate

test:
	@go test ./...

.PHONY: build generate test