package main

import (
	"goapi/config"
	"goapi/pkg/gmail"
	"goapi/pkg/logger"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {
	err := gmail.New().Send("Register Send Code", "123456789", "543619552@qq.com")
	if err != nil {
		logger.Error(err)
	}
}