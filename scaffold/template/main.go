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
	"{{.Vendor}}{{.Service}}/grpc"
	pb "{{.Vendor}}{{.Service}}/proto"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
	syslog "log"
	"net"
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
	srv := grpc.NewRpcServer(logger)

	pb.Register{{title .Service}}Server(srv, new(handler.{{title .Service}}))
	grpc.RegisterHealthServer(srv)

	syslog.Println("serving")
	err = srv.Serve(grpcListener)
	// Run service
	if err != nil {
		logger.Error("during listen err:", zap.Error(err))
	}
}

`,
}
