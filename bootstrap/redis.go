package bootstrap

import "goapi/pkg/redis"

// SetupRedis 初始化Redis
func SetupRedis() {
	redis.ConnectDB()
}

// RedisClose 关闭redis
func RedisClose() {
	redis.Close()
}
