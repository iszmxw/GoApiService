package main

import (
	"fmt"
	"goapi/bootstrap"
	"goapi/config"
	conf "goapi/pkg/config"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {
	// 初始化 SQL
	fmt.Println("初始化 SQL")
	bootstrap.SetupDB()
	// 初始化 Redis
	fmt.Println("初始化 Redis")
	bootstrap.SetupRedis()
	defer bootstrap.RedisClose()
	// 初始化路由绑定
	fmt.Println("加载路由")
	router := bootstrap.SetupRoute()
	// 启动路由
	fmt.Println("启动路由")
	_ = router.Run(fmt.Sprintf(":%s", conf.GetString("app.port")))
}
