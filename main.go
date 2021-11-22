package main

import (
	"fmt"
	"goapi/bootstrap"
	"goapi/config"
	conf "goapi/pkg/config"
	"goapi/pkg/logger"
)

func init() {
	// 初始化配置信息
	config.Initialize()
	// 定义日志目录
	logger.Service = "apiService"
	logger.Init()
}

// @title 用户端接口服务
// @version 3.0
// @description 3.0版本，基于之前的2.0改造的
// @termsOfService http://127.0.0.1/docs/index.html

// @contact.name 追梦小窝
// @contact.url http://github.com/iszmxw
// @contact.email mail@54zm.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1
// @BasePath
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
	if conf.GetString("app.https") == "1" {
		//初始化routes
		_ = router.RunTLS(fmt.Sprintf(":%s", conf.GetString("app.port")), "./config/ssl.pem", "./config/ssl.key")
	} else {
		_ = router.Run(fmt.Sprintf(":%s", conf.GetString("app.port")))
	}
}
