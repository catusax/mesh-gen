package template

// Makefile is the Makefile template used for new projects.
var Makefile = `
GOPATH:=$(shell go env GOPATH)
NAME={{.Service}}
BIN={{.Service}}
REGISTRY_PREFIX={{.RegistryPrefix}}
VERSION={{.Version}}
PORT={{.Port}}
NAMESPACE={{.Namespace}}

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/catusax/mesh-gen/cmd/protoc-gen-go-mesh-gen@latest

.PHONY: proto
proto:
	@cp handler/$(NAME).go proto/$(NAME)_handler.go || true
	@cp handler/$(NAME)_test.go proto/$(NAME)_test.go || true
	@protoc --proto_path=. -I${GOPATH}/src --go-grpc_out=. --go_out=:. --go-mesh-gen_out=. --go-mesh-gen_opt=mesh=traefik,namespace=$(NAMESPACE),port=$(PORT) proto/$(NAME).proto
	@cp proto/$(NAME)_handler.go handler/$(NAME).go
	@cp proto/$(NAME)_test.go handler/$(NAME)_test.go

.PHONY: generate
generate:
	@mesh-gen generate

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -o $(BIN) *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
TAG=
ifeq ($(REGISTRY_PREFIX),)
	TAG=$(NAME)
else
	TAG=$(REGISTRY_PREFIX)/$(NAME)
endif
docker:
	@docker build -t $(TAG):$(VERSION) .

.PHONY: docker-push
docker-push:
	@docker push $(TAG):$(VERSION)

.PHONY: dev
dev:
	@skaffold dev

.PHONY: run
run:
	@skaffold run

`
