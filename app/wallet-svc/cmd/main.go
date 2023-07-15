package main

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"trade-robot-bd/app/wallet-svc/internal/service"
	"trade-robot-bd/app/wallet-svc/server"
	"trade-robot-bd/libs/env"
)

func init() {
}

var (
	id, _ = os.Hostname()
)

func main() {
	log.Println("id:", id)
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{env.EtcdAddr},
	})
	if err != nil {
		log.Fatal(err)
	}
	r := etcd.New(client)
	grpcServers := server.NewGRPCServers(service.NewWalletService())
	httpServer := server.NewHTTPServer()
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(env.WalletSrvName),
		kratos.Version("1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			httpServer,
			grpcServers,
		),
		kratos.Registrar(r),
	)
	//biz.NewWalletRepo().CreateWalletAtRunning()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
