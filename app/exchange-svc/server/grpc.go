package server

import (
	"github.com/go-kratos/kratos/middleware/recovery/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"time"
	v1 "trade-robot-bd/api/exchange/v1"
	"trade-robot-bd/app/exchange-svc/internal/service"
)

// NewGRPCServers new a gRPC server.
func NewGRPCServers(service *service.ExOrderService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			//tracing.Server(
			//	tracing.WithTracerProvider(tp)),
			//logging.Server(logger),
		),
	}
	opts = append(opts, grpc.Timeout(time.Second*5), grpc.Address(":9000"))
	srv := grpc.NewServer(opts...)
	v1.RegisterExOrderServer(srv, service)
	return srv
}
