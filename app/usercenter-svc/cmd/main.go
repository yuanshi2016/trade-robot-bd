package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trade-robot-bd/app/usercenter-svc/internal/service"
	"trade-robot-bd/app/usercenter-svc/server"
	"trade-robot-bd/libs/cache"
	"trade-robot-bd/libs/env"
	"trade-robot-bd/libs/logger"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	id, _ = os.Hostname()
)

func main() {
	/**
	  go run app/common-svc/cmd/main.go
	  		go run app/exchange-svc/cmd/main.go
	  		go run app/grid-strategy-svc/cmd/main.go
	  		go run app/quote-svc/cmd/main.go
	  		go run app/usercenter-svc/cmd/main.go
	  		go run app/wallet-svc/cmd/main.go
	*/
	log.Println("id:", id)
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{env.EtcdAddr},
	})
	if err != nil {
		log.Fatal(err)
	}
	r := etcd.New(client)
	// ctx := context.Background()
	// w, err := r.GetService(ctx, env.WalletSrvName)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var wg sync.WaitGroup
	// wg.Add(1)
	// for _, ve := range w {
	// 	log.Printf("%#v", ve)
	// }
	// log.Println(111)
	// wg.Wait()
	grpcServers := server.NewGRPCServers(service.NewUserService(r))
	httpServer := server.NewHTTPServer()
	defer func() {
		grpcServers.GracefulStop()
		httpServer.Shutdown(context.Background())
	}()
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(env.UserSrvName),
		kratos.Version("1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			httpServer,
			grpcServers,
		),
		kratos.Registrar(r),
	)
	cache.Redis() // 初始化
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func wait() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Info("服务已关闭")
}
