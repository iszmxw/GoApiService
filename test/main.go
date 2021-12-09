package main

import (
	"fmt"
	"goapi/config"
	"goapi/pkg/helpers"
	"goapi/pkg/logger"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {
	// 初始化配置信息
	config.Initialize()
	// 定义日志目录
	logger.Init("test")
}

func main() {
	start := true
	for {
		if start == true {
			for i := 1; i <= 4000; i++ {
				url := "http://10.10.10.80:5000/?page=" + helpers.IntToString(i)
				go func(url string, i int) {
					logger.Info(fmt.Sprintf("执行地%v条数据", i))
					get, err := http.Get(url)
					if err != nil {
						logger.Error(err)
					}
					_, err1 := ioutil.ReadAll(get.Body)
					if err1 != nil {
						logger.Error(err1)
					}
					//logger.Info(string(Body))
				}(url, i)
				if i == 4000 {
					logger.Info("线程任务分配完毕")
				}
				// 每十页数据休息一秒
				if (i % 2) == 0 {
					time.Sleep(time.Second * 1)
				}
			}
			start = false
		}
	}
}
