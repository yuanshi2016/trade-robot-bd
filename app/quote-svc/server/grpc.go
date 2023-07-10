package server

import (
	"github.com/go-kratos/kratos/middleware/recovery/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"time"
	v1 "trade-robot-bd/api/quote/v1"
	"trade-robot-bd/app/quote-svc/internal/service"
)

// NewGRPCServers new a gRPC server.
func NewGRPCServers(service *service.QuoteService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			//tracing.Server(
			//	tracing.WithTracerProvider(tp)),
			//logging.Server(logger),
		),
	}
	opts = append(opts, grpc.Timeout(time.Second*5))
	srv := grpc.NewServer(opts...)
	v1.RegisterQuoteServer(srv, service)
	return srv
}
