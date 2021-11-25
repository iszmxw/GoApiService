package main

import (
	"fmt"
	"github.com/gostudys/huobiapi"
	"goapi/app/models"
	"goapi/app/response"
	"goapi/bootstrap"
	"goapi/config"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/redis"
)

// socket 服务 链接火币网，获取数据缓存到redis

func init() {
	// 初始化配置信息
	config.Initialize()
	// 定义日志目录
	logger.Init("huobiService")
	// 初始化 SQL
	bootstrap.SetupDB()
}

func main() {
	bootstrap.SetupRedis("15")
	defer bootstrap.RedisClose()
	var KLineCode []response.KLineCode
	mysql.DB.Model(models.Currency{}).Find(&KLineCode)
	type topics struct {
		Topic        string `json:"topic"`         // K线图代码
		DecimalScale int    `json:"decimal_scale"` // 自有币位数
	}
	var Topics []topics
	for _, val := range KLineCode {
		// 收集所有要订阅处理的 topic
		Topics = append(Topics,
			// 1分钟k线图
			topics{
				Topic:        "market." + val.KLineCode + ".kline.1min",
				DecimalScale: val.DecimalScale,
			},
			// 5分钟k线图
			topics{
				Topic:        "market." + val.KLineCode + ".kline.5min",
				DecimalScale: val.DecimalScale,
			},
			// 行情
			topics{
				Topic:        "market." + val.KLineCode + ".depth.step1",
				DecimalScale: val.DecimalScale,
			},
			// K线图历史 1分钟
			topics{
				Topic:        "market." + val.KLineCode + ".kline.1min",
				DecimalScale: val.DecimalScale,
			},
			// K线图历史 5分钟
			topics{
				Topic:        "market." + val.KLineCode + ".kline.5min",
				DecimalScale: val.DecimalScale,
			},
			// K线图历史 30分钟
			topics{
				Topic:        "market." + val.KLineCode + ".kline.30min",
				DecimalScale: val.DecimalScale,
			},
			// K线图历史 60分钟
			topics{
				Topic:        "market." + val.KLineCode + ".kline.60min",
				DecimalScale: val.DecimalScale,
			},
			// K线图历史 60分钟
			topics{
				Topic:        "market." + val.KLineCode + ".kline.1day",
				DecimalScale: val.DecimalScale,
			})
	}
	start := "ok"
	for {
		if start == "ok" {
			logger.Info(Topics)
			for _, obj := range Topics {
				logger.Info("缓存 " + obj.Topic)
				go SocketHuoBi(obj.Topic)
			}
			start = "no"
		}
	}
}

// 链接火币网

func SocketHuoBi(parentTopic string) {
	market, NewMarketErr := huobiapi.NewMarket()
	if NewMarketErr != nil {
		logger.Error(NewMarketErr)
		return
	}
	SubscribeErr := market.Subscribe(parentTopic, func(topic string, json *huobiapi.JSON) {
		logger.Info(fmt.Sprintf("成功订阅火币网，您当前订阅的 topic 为：%v", parentTopic))
		// 收到数据更新时回调
		logger.Info(topic)
		logger.Info(json)
		// 火币网推送的回来的数据转换为字符串
		msgData, MarshalJSONErr := json.MarshalJSON()
		if MarshalJSONErr != nil {
			logger.Error(MarshalJSONErr)
			return
		}
		if len(msgData) > 0 {
			// 收集数据缓存到redis
			_, redisErr := redis.Add(topic, msgData, 0)
			if redisErr != nil {
				logger.Error(redisErr)
				return
			}
			// 收集数据缓存到redis
		}
	})
	// socket 火币网失败，释放资源
	if SubscribeErr != nil {
		logger.Error(SubscribeErr)
		CloseErr := market.Close()
		if CloseErr != nil {
			logger.Error(CloseErr)
			return
		}
	}
}
