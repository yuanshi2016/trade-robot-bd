package exchangeclient

import (
	"net/http"
	"net/url"
	"strings"
	"time"
	"trade-robot-bd/libs/env"
	"trade-robot-bd/libs/goex"
)

type ExClientI interface {
	GetAccountSpot() (*goex.Account, error)
	GetAccountSwap() (*goex.Account, error)
	CheckIfApiValid() error
}

// todo 代理
var (
	client = &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   strings.Split(env.ProxyAddr, "//")[1],
				}, nil
			},
		},
	}
)
