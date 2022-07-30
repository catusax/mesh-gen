package template

// Dockerfile is the Dockerfile template used for new projects.
var Dockerfile = Template{
	Path: "Dockerfile",
	Value: `FROM golang:alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc
WORKDIR /go/src/{{.Service}}
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,mode=0777,id=gomod,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,mode=0777,target=/root/.cache/go-build \
    --mount=type=cache,mode=0777,id=gomod,target=/go/pkg/mod \
    go mod tidy && go build -o {{.Service}}

FROM ghcr.io/catusax/grpc-runner:latest
ENV CONTAINER=docker
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/{{.Service}}/{{.Service}} /{{.Service}}
ENTRYPOINT ["/{{.Service}}"]
CMD []
`,
}

// DockerIgnore is the .dockerignore template used for new projects.
var DockerIgnore = Template{
	Path: ".dockerignore",
	Value: `.gitignore
Dockerfile
resources/
skaffold.yaml
`,
}
