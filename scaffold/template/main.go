package template

type Template struct {
	Path  string
	Value string
}

// MainSRV is the main template used for new service projects.
var MainSRV = Template{
	Path: "main.go",
	Value: `package main

import (
	"{{.Vendor}}{{.Service}}/handler"
	pb "{{.Vendor}}{{.Service}}/proto"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
	syslog "log"
	"net"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
    "google.golang.org/grpc/xds"
)

var (
	service = "{{lower .Service}}.service"
)


func main() {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zapcore.WarnLevel))
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger.Named(service))

	grpcListener, err := net.Listen("tcp", ":{{.Port}}")
	// Create service
	srv := NewRpcServer(logger)

	pb.Register{{title .Service}}Server(srv, new(handler.{{title .Service}}))

	syslog.Println("serving")
	err = srv.Serve(grpcListener)
	// Run service
	if err != nil {
		logger.Error("during listen err:", zap.Error(err))
	}
}

type GRPCServer interface {
	grpc.ServiceRegistrar
	Serve(lis net.Listener) error
}

func NewRpcServer(logger *zap.Logger) GRPCServer {
	options := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zap.L()),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zap.L()),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}

	if os.Getenv("CONTAINER") != "" {
		return xds.NewGRPCServer(options...)
	} else {
		return grpc.NewServer(options...)
	}
}
`,
}
