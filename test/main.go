package main

import (
	"fmt"
	"goapi/config"
	conf "goapi/pkg/config"
	"goapi/pkg/email/qq"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {
	err := qq.SendEmail("code", conf.GetString("email.qq.user"), "543619552@qq.com", "121212")
	if err != nil {
		fmt.Println(err)
	}
}
