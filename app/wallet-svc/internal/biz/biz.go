package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	pb "trade-robot-bd/api/usercenter/v1"
	"trade-robot-bd/app/wallet-svc/cache"
	"trade-robot-bd/app/wallet-svc/internal/dao"
	"trade-robot-bd/libs/env"
	"trade-robot-bd/libs/exchangeclient"
)

const (
	ErrID = "wallet"
	//BinanceApiKey = "fev72IlrChwPbO8Yp3D57RkvIiuUwkIFK3dJoQi7cQaYyv00DiBwxDiXm4DH4HZq"
	//BinanceSecret = "YGvVhns0OlIxMJ1of4apa0IeYGbXsFvCrbewrTYveQz0qfxDhRalBfBJd7EUN4iP"
	BinanceApiKey = "lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d"
	BinanceSecret = "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz"
)

type WalletRepo struct {
	dao          *dao.Dao
	cacheService *cache.Service
	binance      *exchangeclient.Binance
	UserSrv      pb.UserClient
}

func NewWalletRepo() *WalletRepo {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{env.EtcdAddr},
	})
	if err != nil {
		log.Fatal(err)
	}
	r := etcd.New(client)
	ctx := context.Background()
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%v", env.UserSrvName)),
		grpc.WithDiscovery(r),
	)
	return &WalletRepo{
		dao:          dao.New(),
		cacheService: cache.NewService(),
		binance:      exchangeclient.InitBinance(BinanceApiKey, BinanceSecret),
		UserSrv:      pb.NewUserClient(conn),
	}
}
