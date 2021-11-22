package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goapi/app/response"
	"goapi/pkg/huobi"
	"goapi/pkg/logger"
	"goapi/pkg/redis"
	"goapi/pkg/websocket"
	"net/http"
	"strings"
	"time"
)

// k线图服务

type KlineController struct {
	BaseController
}

// WsHandler socket 负责转发redis缓存的socket数据
func (h *KlineController) WsHandler(c *gin.Context) {
	RequestId := logger.RequestId
	ws, err := wss.GetSocket(c)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("客户端建立socket成功")
	// 为了释放for死循环的资源
	for {
		// 读取用户客户端ws中订阅的 topic
		mt, message, ReadMessageErr := ws.Conn.ReadMessage()
		if ReadMessageErr != nil {
			marketErr := ws.Close()
			if marketErr != nil {
				logger.Error(ReadMessageErr)
				logger.Error(marketErr)
			} else {
				logger.Info("关闭成功")
			}
			// 跳出for关闭后面的读取
			break
		}
		//对数据进行切割，读取参数
		//如果请求的是 "market.btcusdt.kline.1min" ,订阅这条信息，然后再返回
		msg := string(message)
		msg = strings.Trim(msg, "\"")
		logger.Info(msg)

		// 24 小时内一个长连接禁止重复订阅相同的 topic
		checkMsg := RequestId + ":" + msg
		if redis.CheckExist(checkMsg) {
			logger.Info("该订阅已存在，请勿重新订阅")
			continue
		}
		_, redisAddErr := redis.Add(checkMsg, msg, 60*60*24)
		if redisAddErr != nil {
			logger.Error(redisAddErr)
			return
		}
		// 24 小时内一个长连接禁止重复订阅相同的 topic

		// 开启协诚处理推送消息
		go func(ws *wss.WsConn, msg string) {
			// 为了释放for死循环的资源
			for {
				push, getErr := redis.Get(msg)
				if getErr != nil {
					logger.Error(getErr)
					return
				}
				// 空消息跳过
				if len(push) <= 0 {
					continue
				}
				// 1秒推送一次消息
				time.Sleep(time.Second * 1)
				ws.Mux.Lock()
				WriteMessageErr := ws.Conn.WriteMessage(mt, []byte(push))
				ws.Mux.Unlock()
				if WriteMessageErr != nil {
					// 跳出for关闭后面的推送
					break
				}
				// 返回数据给用户
			}
			logger.Info("已经跳出协诚的循环推送")
		}(ws, msg)
	}
	logger.Info("socket被关闭，本次链接结束")
}

// HistoryHandler 历史行情
func (h *KlineController) HistoryHandler(c *gin.Context) {
	symbol := c.Query("symbol")
	period := c.Query("period")
	if len(symbol) <= 0 || len(period) <= 0 {
		symbol = "btcusdt"
		period = "1min"
	}
	kline, err := huobi.RequestHuoBiKline(symbol, period)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{
			//返回数据
			"data": nil,
		})
		return
	}
	var data response.KlineData
	err = json.Unmarshal(kline, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		//返回数据
		"data": data,
	})
}
