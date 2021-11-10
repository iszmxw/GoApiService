package huobi

import (
	"errors"
	"fmt"
	"goapi/pkg/request"
)

type KData struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Data   []struct {
		Id     int     `json:"id"`
		Open   float64 `json:"open"`
		Close  float64 `json:"close"` // 获取本阶段收盘价
		Low    float64 `json:"low"`   // 买入的时候取 low
		High   float64 `json:"high"`  // 卖出的时候取 high
		Amount float64 `json:"amount"`
		Vol    float64 `json:"vol"`
		Count  int     `json:"count"`
	} `json:"data"`
}

//Close  float64 `json:"close"` // 获取本阶段收盘价
//Low    float64 `json:"low"`   // 买入的时候取 low
//High   float64 `json:"high"`  // 卖出的时候取 high

func Kline(kline, name string) (float64, error) {
	var Price interface{}
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
	if name == "close" {
		// 获取本阶段收盘价
		Price = resp["data"].([]interface{})[0].(map[string]interface{})["close"]
		//fmt.Println(fmt.Sprintf("%T", Price))
	}
	if name == "low" {
		// 买入的时候取 low
		Price = resp["data"].([]interface{})[0].(map[string]interface{})["low"]
	}
	if name == "high" {
		// 卖出的时候取 high
		Price = resp["data"].([]interface{})[0].(map[string]interface{})["high"]
	}
	return Price.(float64), nil
}
