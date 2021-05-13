SHELL := /bin/bash

build: generate
	@go build -o bin/service

generate:
	@buf generate

test:
	@go test ./...

.PHONY: build generate test