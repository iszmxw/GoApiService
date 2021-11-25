package bootstrap

import "goapi/pkg/redis"

// SetupRedis 初始化Redis
func SetupRedis(selectDB string) {
	redis.ConnectDB(selectDB)
}

// RedisClose 关闭redis
func RedisClose() {
	redis.Close()
}
