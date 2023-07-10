package main

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"

	"github.com/go-kratos/kratos/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"trade-robot-bd/app/common-svc/internal/service"
	"trade-robot-bd/app/common-svc/server"
	"trade-robot-bd/libs/env"
)

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
	grpcServers := server.NewGRPCServers(service.NewCommonService())
	httpServer := server.NewHTTPServer()
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(env.CommonSrvName),
		kratos.Version("1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			httpServer,
			grpcServers,
		),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
