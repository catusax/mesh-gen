package template

// MainSRV is the main template used for new service projects.
var MainSRV = `package main

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


func NewRpcServer(logger *zap.Logger) *grpc.Server {
	return grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zap.L()),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zap.L()),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
}
`
