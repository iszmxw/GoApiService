package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"goapi/pkg/config"
	"time"
)

// Redis redis.Client 对象
var Redis *redis.Client
var expireTime = 600

// ConnectDB 初始化模型
func ConnectDB() {
	// 初始化 Redis 连接信息
	var (
		err       error
		RedisIp   = config.GetString("redis.host")
		RedisPort = config.GetString("redis.port")
		DefaultDB = config.GetInt("redis.db")
		Pw        = config.GetString("redis.password")
	)
	if len(Pw) > 0 {
		Redis = redis.NewClient(&redis.Options{
			Addr:     RedisIp + ":" + RedisPort,
			DB:       DefaultDB, // use default DB
			Password: Pw,        // no password set
		})
	} else {
		Redis = redis.NewClient(&redis.Options{
			Addr: RedisIp + ":" + RedisPort,
			DB:   DefaultDB, // use default DB
		})
	}
	_, err = Redis.Ping().Result()
	if err != nil {
		fmt.Println("redis连接错误" + err.Error())
	}
}

func CheckExist(key string) bool {
	a, err := Redis.Exists(key).Result()
	if err != nil {
		fmt.Println("判断key存在失败")
		return false
	}
	if a == 1 {
		fmt.Println("key存在")
		return true
	}
	return false
}

func Add(key string, value interface{}, exTime int) (bool, error) {
	if exTime >= 0 {
		expireTime = exTime
	}
	err := Redis.Set(key, value, time.Duration(expireTime)*time.Second).Err()
	if err != nil {
		fmt.Println("设置key失败")
		return false, err
	}
	return true, nil
}

func Delete(key string) bool {
	err := Redis.Del(key).Err()
	if err != nil {
		fmt.Println("删除key失败" + err.Error())
		return false
	}
	return true
}

func Get(key string) (string, error) {
	value, err := Redis.Get(key).Result()
	return value, err
}

// Redis Keys 命令 - 查找所有符合给定模式( pattern)的 key
// https://www.redis.net.cn/order/3535.html

func Keys(key string) ([]string, error) {
	value, err := Redis.Keys(key).Result()
	return value, err
}

func SubExpireEvent(channels string) *redis.PubSub {
	// 订阅key过期事件
	//sub := Redis.Subscribe("__keyevent@0__:expired")
	return Redis.Subscribe(channels)
	// 这里通过一个for循环监听redis-server发来的消息。
	// 当客户端接收到redis-server发送的事件通知时，
	// 客户端会通过一个channel告知我们。我们再根据
	// msg的channel字段来判断是不是我们期望收到的消息，
	// 然后再进行业务处理。
	//for {
	//	msg := <-sub.Channel()
	//	fmt.Println("Channel ", msg.Channel)
	//	fmt.Println("pattern ", msg.Pattern)
	//	fmt.Println("pattern ", msg.Payload)
	//}
}

func Close() {
	err := Redis.Close()
	if err != nil {
		fmt.Println("RedisCloseError", err.Error())
	}
}
