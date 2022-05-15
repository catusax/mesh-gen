# mesh-gen

microservices have a lot of duplicate basic codes, you have to copy-paste them for every service.
mesher can generate that codes for you.

- generate project scaffold, kubernetes resources, and dockerfile
- generate testing, handler, and client code form protobuf definition file

## mesh-gen cli

### TL;DR:

just run this ,then your new microservice is ready for business logic

```shell
go install github.com/catusax/mesh-gen/cmd/mesh-gen@latest
mesh-gen new my-service
cd my-service

# edit proto/my-service.proto

make init proto tidy
```

### config

mesher reads config form `Makefile`,you can edit variables in `Makefile`, then regenerate k8s resource.

```shell
mesh-gen generate
make proto
```

## protoc-gen-go-mesh-gen

generate testing , handler , and client code from proto file

demo: [demo](cmd/protoc-gen-go-mesh-gen/proto)

### usage

```shell
## dep
go get -u google.golang.org/protobuf/proto
go install github.com/golang/protobuf/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

go install github.com/catusax/mesh-gen/cmd/protoc-gen-mesh-genc@latest
protoc --proto_path=. --go-grpc_out=. --go_out=:. --go-mesh-gen=. --go-mesh-gen_opt=mesh=<traefik or istio>,namespace=<k8s namespace to put your service>,port=<service port for client connecting> proto/<file>.proto

```

## thanks

<https://github.com/rotemtam/protoc-gen-go-ascii>

<https://github.com/asim/go-micro>

## License

Copyright (c) 2022 catusax

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
