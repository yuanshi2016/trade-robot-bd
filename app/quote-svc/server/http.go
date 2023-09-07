package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/middleware/recovery/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"time"
	"trade-robot-bd/app/quote-svc/router"
	"trade-robot-bd/libs/logger"
	"trade-robot-bd/libs/middleware"
)

const (
	port = ":9530"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer() *http.Server {
	engine := gin.Default()
	//engine.Use(middleware.CorsL())
	engine.Use(middleware.CorsR())
	engine.Use(kgin.Middlewares(recovery.Recovery()))
	httpSrv := http.NewServer(
		http.Address(port),
		http.Timeout(time.Second*10),
	)
	router.Init(engine)
	httpSrv.HandlePrefix("/", engine)
	pprof.Register(engine, "/quote/debug")
	logger.Infof("启动服务，监听端口：%v", port)
	return httpSrv
}
