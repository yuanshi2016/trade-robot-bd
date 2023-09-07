package env

var (
	developEnv = envConfig{
		runMode:   "dev",
		redisAddr: "103.158.36.177:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "trade-etcd.develop:2379",
		//etcdAddr:  "103.158.36.177:12379",
		dbDSN:     "wq_fotune:Yuanshi20188@tcp(103.158.36.177:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade:Yuanshi20188@103.158.36.177:38888/ifortune",
	}
	releaseEnv = envConfig{
		runMode:   "dev",
		redisAddr: "103.158.36.177:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "trade-etcd.develop:2379",
		dbDSN:     "wq_fotune:Yuanshi20188@tcp(103.158.36.177:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade:Yuanshi20188@103.158.36.177:38888/ifortune",
	}
	proEnv = envConfig{
		runMode:   "dev",
		redisAddr: "103.158.36.177:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "trade-etcd.develop:2379",
		dbDSN:     "wq_fotune:Yuanshi20188@tcp(103.158.36.177:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade:Yuanshi20188@103.158.36.177:38888/ifortune",
	}
	testEnv = envConfig{
		runMode:   "dev",
		redisAddr: "103.158.36.177:6379",
		redisPWD:  "Yuanshi20188",
		etcdAddr:  "trade-etcd.develop:2379",
		dbDSN:     "wq_fotune:Yuanshi20188@tcp(103.158.36.177:3306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://trade:Yuanshi20188@103.158.36.177:38888/ifortune",
	}
)
