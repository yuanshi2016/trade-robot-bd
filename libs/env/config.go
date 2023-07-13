package env

import (
	"net/http"
	"os"
)

const (
	runMode   = "RUN_MODE"
	redisAddr = "REDIS_ADDR"
	redisPWD  = "REDIS_PWD"
	dbDSN     = "DB_DSN"
	etcdAddr  = "ETCD_ADDR"
	mongoAddr = "MONGO_ADDR"

	mode = "MODE"
)

var (
	configMap = initial()

	RunMode   = configMap.getValue(runMode)
	RedisAddr = configMap.getValue(redisAddr)
	RedisPWD  = configMap.getValue(redisPWD)
	DbDSN     = configMap.getValue(dbDSN)
	EtcdAddr  = configMap.getValue(etcdAddr)
	MongoAddr = configMap.getValue(mongoAddr)

	GridNum = 10
	// MsgCenterURL 上报异常到消息中心URL
	MsgCenterURL = "http://localhost:20080/v1/dataCollect/systemMsg"
	// ExchangeAccessURL 获取交易所访问授权信息URL，url后面需要加参数 /:userid/:apikey
	ExchangeAccessURL = "https://yun.local100.com/test/exchange/v1/exchange/apiInfo"
	// StatisticalInfoURL 获取统计信息URL url后面需要加参数 /:user_id/:strategyId
	StatisticalInfoURL = "https://yun.local100.com/test/exchange/v1/user/strategy/evaluationNoAuth"
	// NotifyStatisticsURL 通知统计URL
	NotifyStatisticsURL = "https://yun.local100.com/test/exchange/v1/forward-offer/orderGrid"
	// NotifyStrategyStartUpURL 启动策略通知接口
	NotifyStrategyStartUpURL = "https://yun.local100.com/test/wallet/v1/wallet/strategyStartUpNotify"
	//ProxyAddr                = ""
	ProxyAddr = "socks5://10.10.1.3:10801"

	ExchangeSrvName = "exchange-order.srv"
	UserSrvName     = "usercenter.srv"
	WalletSrvName   = "wallet.srv"
	QuoteSrvName    = "quote.srv"
	CommonSrvName   = "common.srv"
)

type envConfig map[string]string

func initial() envConfig {
	var config envConfig
	releaseMode := os.Getenv(mode)
	if releaseMode == "production" {
		config = proEnv
	} else if releaseMode == "release" {
		config = releaseEnv
	} else {
		config = developEnv
	}
	return config
}

func (env envConfig) getValue(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if v, ok := env[key]; ok {
		return v
	}
	return ""
}

func GetProxyHttpClient() *http.Client {
	client := &http.Client{}
	//client.Transport = &http.Transport{
	//	Proxy: func(req *http.Request) (*url.URL, error) {
	//		return &url.URL{
	//			Scheme: "socks5",
	//			Host:   strings.Split(ProxyAddr, "//")[1],
	//		}, nil
	//	},
	//}
	return client
}
