package server

import (
	"time"
	"trade-robot-bd/app/wallet-svc/router"
	"trade-robot-bd/libs/logger"
	"trade-robot-bd/libs/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/middleware/recovery/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
)

const (
	port = ":9531"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer() *http.Server {
	engine := gin.Default()
	engine.Use(kgin.Middlewares(recovery.Recovery()))
	router.Init(engine)
	httpSrv := http.NewServer(http.Address(port), http.Timeout(time.Second*10), middleware.KoCors())
	httpSrv.HandlePrefix("/", engine)
	pprof.Register(engine, "/wallet/debug")
	logger.Infof("启动服务，监听端口：%v", port)
	return httpSrv
}
