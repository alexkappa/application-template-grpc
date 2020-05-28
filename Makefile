SHELL := /bin/bash
PROTO_FILES := $(shell ls api/proto/*.proto api/*/proto/*.proto)

build: generate
	@go build -o bin/service

generate:
	@protoc -I . \
		-I third_party/googleapis \
		-I third_party/grpc-gateway \
		--go_out=plugins=grpc,paths=source_relative:. \
		--grpc-gateway_out=logtostderr=true:. \
		--swagger_out=logtostderr=true,allow_merge=true,use_go_templates=true,merge_file_name=docs:. \
		$(PROTO_FILES)
	@go generate ./...

test:
	@go test ./...

.PHONY: build generate test