package template

// Makefile is the Makefile template used for new projects.
var Makefile = Template{
	Path: "Makefile",
	Value: `
GOPATH:=$(shell go env GOPATH)
BIN={{.Service}}
NAME={{.Service}}
REGISTRY_PREFIX={{.RegistryPrefix}}
VERSION={{.Version}}
PORT={{.Port}}
NAMESPACE={{.Namespace}}
MESH={{.Mesh}}

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/catusax/mesh-gen/cmd/protoc-gen-go-mesh-gen@latest
	@go install github.com/catusax/mesh-gen/cmd/mesh-gen@latest

.PHONY: proto
proto:
	@protoc --proto_path=. -I${GOPATH}/src --go-grpc_out=. --go_out=:. --go-mesh-gen_out=. --go-mesh-gen_opt=mesh==$(MESH),handler=handler,namespace=$(NAMESPACE),port=$(PORT),name=$(NAME) proto/$(NAME).proto

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
	@CGO_ENABLED=0 GOOS=linux go build -o $(BIN) *.go

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
	@DOCKER_BUILDKIT=1 docker build -t $(TAG):$(VERSION) .

.PHONY: docker-push
docker-push:
	@docker push $(TAG):$(VERSION)

.PHONY: dev
dev:
	@skaffold dev

.PHONY: run
run:
	@skaffold run --tail

`,
}
