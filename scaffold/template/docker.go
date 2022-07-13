package template

// Dockerfile is the Dockerfile template used for new projects.
var Dockerfile = `FROM golang:alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc
WORKDIR /go/src/{{.Service}}
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    make tidy build

FROM scratch
ENV CONTAINER=docker
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/{{.Service}}/{{.Service}} /{{.Service}}
ENTRYPOINT ["/{{.Service}}"]
CMD []
`

// DockerIgnore is the .dockerignore template used for new projects.
var DockerIgnore = `.gitignore
Dockerfile
resources/
skaffold.yaml
`
