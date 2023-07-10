package env

var (
	developEnv = envConfig{
		runMode:   "dev",
		redisAddr: "10.10.1.100:6379",
		redisPWD:  "Yuanshi20188",
		//etcdAddr:  "192.168.5.5:12379",
		etcdAddr:  "10.10.1.10:2379",
		dbDSN:     "root:root@tcp(127.0.0.1:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade:199535@10.10.1.100:38888/ifortune",
	}
	releaseEnv = envConfig{
		runMode:   "dev",
		redisAddr: "10.10.1.100:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "10.10.1.10:2379",
		dbDSN:     "root:root@tcp(127.0.0.1:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade:199535@10.10.1.100:38888/ifortune",
	}
	proEnv = envConfig{
		runMode:   "dev",
		redisAddr: "10.10.1.100:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "10.10.1.10:2379",
		dbDSN:     "root:root@tcp(127.0.0.1:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade:199535@10.10.1.100:38888/ifortune",
	}
)
