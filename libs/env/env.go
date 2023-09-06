package env

var (
	developEnv = envConfig{
		runMode:   "dev",
		redisAddr: "10.10.1.100:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "10.10.1.100:2379",
		dbDSN:     "root:root@tcp(10.10.1.100:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade1:199535@10.10.1.100:38888/ifortune",
	}
	releaseEnv = envConfig{
		runMode:   "dev",
		redisAddr: "10.10.1.100:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "10.10.1.100:2379",
		dbDSN:     "root:root@tcp(10.10.1.100:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade1:199535@10.10.1.100:38888/ifortune",
	}
	proEnv = envConfig{
		runMode:   "dev",
		redisAddr: "10.10.1.100:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "10.10.1.100:2379",
		dbDSN:     "root:root@tcp(10.10.1.100:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade1:199535@10.10.1.100:38888/ifortune",
	}
	testEnv = envConfig{
		runMode:   "dev",
		redisAddr: "10.10.1.100:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "10.10.1.6:2379",
		dbDSN:     "root:root@tcp(10.10.1.100:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade1:199535@10.10.1.100:38888/ifortune",
	}
)
