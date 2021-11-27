package bootstrap

import (
	"goapi/pkg/redis"
	"goapi/pkg/redis_socket"
)

// SetupRedis 初始化Redis
func SetupRedis(selectDB string) {
	redis.ConnectDB(selectDB)
}

// SetupSocketRedis 初始化Redis
func SetupSocketRedis(selectDB string) {
	redis_socket.ConnectDB(selectDB)
}

// RedisClose 关闭redis
func RedisClose() {
	redis.Close()
	redis_socket.Close()
}
