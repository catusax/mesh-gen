GOPATH:=$(shell go env GOPATH)
NAME=wallet

.PHONY: default
default: build

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go get -u github.com/golang/protobuf

.PHONY: proto
proto:
	@protoc --proto_path=. -I${GOPATH}/src --go_out=. --go-grpc_out=. --go-kit_out=.  test.proto

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go get -u protoc-gen-go-kit
	@go install .

