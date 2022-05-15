package template

// Dockerfile is the Dockerfile template used for new projects.
var Dockerfile = `FROM golang:alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/{{.Service}}
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc
COPY Makefile go.mod go.sum ./
RUN go mod download
COPY . .
RUN make tidy build

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
