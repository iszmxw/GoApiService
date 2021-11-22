package main

import (
	"encoding/base64"
	"fmt"
	"goapi/app/models"
	"goapi/app/response"
	"goapi/bootstrap"
	"goapi/config"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"io/ioutil"
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
	//db.Debug().Model(models.User{}).Where("id", "503").Find(&user)
	var v1 models.Verify
	db.Debug().Model(models.Verify{}).Where("user_id", "501").Find(&v1)
	imgbase64 := v1.ImgCardFront
	imgbff, _ := base64.StdEncoding.DecodeString(imgbase64)
	fmt.Println(len(imgbff))
	err2 := ioutil.WriteFile("./test.jpg", imgbff, 0666)
	if err2 != nil {
		logger.Error(err2)
		return
	}
	fmt.Println(user)
}
