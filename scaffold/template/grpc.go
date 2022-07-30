package template

import "path/filepath"

var GRPCMiddleWare = Template{
	Path: filepath.Join("grpc", "middleware.go"),
	Value: `package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/xds"
	"net"
	"os"
)

type Server interface {
	grpc.ServiceRegistrar
	Serve(lis net.Listener) error
}

func NewRpcServer(logger *zap.Logger) Server {
	options := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
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

var GRPCHealth = Template{
	Path: filepath.Join("grpc", "health.go"),
	Value: `package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
)

var _health *health.Server
var _healthOnce sync.Once

func GetHealthServer() grpc_health_v1.HealthServer {
	_healthOnce.Do(func() {
		_health = health.NewServer()
	})
	return _health
}

func RegisterHealthServer(s grpc.ServiceRegistrar) {
	grpc_health_v1.RegisterHealthServer(s, GetHealthServer())
}

`,
}
