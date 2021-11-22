package main

import (
	"encoding/base64"
	"fmt"
	"goapi/app/models"
	"goapi/bootstrap"
	"goapi/config"
	"goapi/pkg/mysql"
	"io/ioutil"
	"os"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}
func base64transfer(base64str string, filedir string, filename string) error {
	//生成文件路径
	Eerr := os.Mkdir(filedir, 0755)
	if Eerr != nil {
		//logger.Error(Eerr)
		return Eerr
	}
	filepath := filedir + filename
	println(filepath)
	imgbff, _ := base64.StdEncoding.DecodeString(base64str)
	if len(imgbff) > 5120000 {
		//return errors.New("每张图片大小不能超过5mb")
	}
	err := ioutil.WriteFile(filepath, imgbff, 0666)
	if err != nil {
		//logger.Error(err)
		return err
	}
	return nil
}

func main() {
	// 初始化 SQL
	fmt.Println("初始化 SQL")
	bootstrap.SetupDB()
	//var user response.User
	db := mysql.DB
	//db.Debug().Model(models.User{}).Where("id", "503").Find(&user)
	var v1 models.Verify
	db.Debug().Model(models.Verify{}).Where("user_id", "501").Find(&v1)
	firdir := "./resource/photo/555/"
	filename := "./test.jpg"
	imgbase64 := v1.ImgCardFront
	err := base64transfer(imgbase64, firdir, filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("1111111")
	//imgbff, _ := base64.StdEncoding.DecodeString(imgbase64)
	//fmt.Println(len(imgbff))
	//len(imgbff)
	//c.Writer.WriteString(string(imgbff))这样应该就可以了
	//logger.Info(v1)
	//if err != nil {
	//	logger.Error(err)
	//}
	//filepath := fmt.Sprintf("./resource/phone/%d/", v1.UserId)
	//filename := fmt.Sprintf("./resource/phone/%d/%s.jpg", v1.UserId, "cf")
	//_ = ioutil.WriteFile("./test.jpg", imgbff, 0666)
	//file, _ := ioutil.ReadFile(filename)
	//c.Writer.WriteString(string(file))
	//if err2 != nil {
	//	logger.Error(err2)
	//	return
	//}

	//fmt.Println(user)
}
