GOPATH:=$(shell go env GOPATH)
NAME=wallet

.PHONY: default
default: build

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go get -u github.com/golang/protobuf

.PHONY: proto
proto:
	@protoc --proto_path=. -I${GOPATH}/src --go_out=. --go-grpc_out=. --go-mesher_out=. --go-mesher_opt=test=true,handler=true test.proto

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go get -u github.com/catusax/go-kit-micro/cmd/protoc-gen-go-kit-grpc
	@go install .

