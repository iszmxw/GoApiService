package main

import (
	"fmt"
	"goapi/app/models"
	"goapi/app/response"
	"goapi/bootstrap"
	"goapi/config"
	"goapi/pkg/mysql"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {
	// 初始化 SQL
	fmt.Println("初始化 SQL")
	bootstrap.SetupDB()
	var user response.User
	db := mysql.DB
	db.Debug().Model(models.User{}).Where("id", "503").Find(&user)
	fmt.Println(user)
}
