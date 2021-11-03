package huobi

import (
	"errors"
	"fmt"
	"goapi/pkg/request"
)

func Kline(kline string) (float64, error) {
	url := "https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=" + kline
	resp, err := request.Get(url)
	if err != nil {
		fmt.Println("GET请求错误", err.Error())
		return 0, err
	}
	if resp["status"] != "ok" {
		// todo 重新记录回redis
		msg := fmt.Sprintf("请求火币网获取数据失败，请求url：%v，status：%v", url, resp["status"])
		fmt.Println(msg)
		return 0, errors.New(msg)
	}
	fmt.Println(resp)
	// 获取本阶段收盘价
	clinchPrice := resp["data"].([]interface{})[0].(map[string]interface{})["close"]
	//fmt.Println(fmt.Sprintf("%T", clinchPrice))
	return clinchPrice.(float64), nil
}
