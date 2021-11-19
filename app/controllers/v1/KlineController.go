package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gostudys/huobiapi"
	"goapi/pkg/logger"
	"goapi/pkg/websocket"
	"strings"
)

// k线图服务

type KlineController struct {
	BaseController
}

// WsHandler socket
func (h *KlineController) WsHandler(c *gin.Context) {
	ws, err := wss.GetSocket(c)
	if err != nil {
		logger.Error(err)
		return
	}
	for {
		market, NewMarketErr := huobiapi.NewMarket()
		if NewMarketErr != nil {
			logger.Error(NewMarketErr)
			continue
		}
		//读取ws中的数据
		mt, message, ReadMessageErr := ws.Conn.ReadMessage()
		if ReadMessageErr != nil {
			logger.Error(ReadMessageErr)
			marketErr := market.Close()
			if marketErr != nil {
				logger.Error(marketErr)
				continue
			} else {
				logger.Info("关闭成功")
			}
			continue
		}
		logger.Info(fmt.Sprintf("当前数据类型为: %v", mt))
		//对数据进行切割，读取参数
		//如果请求的是market.ethbtc.kline.5min,订阅这条信息，然后再返回
		msg := string(message)
		newMsg := string([]byte(msg)[1 : len([]byte(msg))-1])
		//打印请求参数
		logger.Info(newMsg)
		if strings.Contains(msg, "1min") || strings.Contains(msg, "step1") {
			logger.Info("当前为1分钟请求")
			SubscribeErr := market.Subscribe(newMsg, func(topic string, json *huobiapi.JSON) {
				logger.Info(msg)
				logger.Info(newMsg)
				if err != nil {
					logger.Error(err)
				}
				logger.Info("订阅成功")
				// 收到数据更新时回调
				logger.Info(topic)
				logger.Info(json)
				// 火币网推送的回来的数据转换为字符串
				msgData, MarshalJSONErr := json.MarshalJSON()
				if MarshalJSONErr != nil {
					logger.Error(MarshalJSONErr)
					return
				}
				ws.Mux.Lock()
				WriteMessageErr := ws.Conn.WriteMessage(mt, msgData)
				ws.Mux.Unlock()
				if WriteMessageErr != nil {
					logger.Error(WriteMessageErr)
					wsErr := ws.Close()
					if wsErr != nil {
						logger.Error(wsErr)
					}
					return
				}
				//返回数据给用户
			})
			// socket 火币网失败，释放资源
			if SubscribeErr != nil {
				CloseErr := market.Close()
				if CloseErr != nil {
					logger.Error(CloseErr)
					continue
				}
			}
		} else {
			//写入ws数据
			//go func() {
			//	for {
			//
			//		go func() {
			//			err = market.Subscribe(newMsg, func(topic string, hjson *huobiapi.JSON) {
			//				//logger.Info(msg)
			//				if err != nil {
			//					logger.Error(err)
			//				}
			//				//订阅成功
			//				//logger.Info("订阅成功")
			//				//120后自动取消订阅
			//				go func() {
			//					time.Sleep(60 * time.Minute)
			//					//logger.Info("取消订阅成功")
			//					market.Unsubscribe(newMsg)
			//					//market.ReceiveTimeout
			//
			//				}()
			//
			//				// 收到数据更新时回调
			//				//logger.Info(topic)
			//				//logger.Info(hjson)
			//				jsondata, MarshalJSONErr := hjson.MarshalJSON()
			//				if err != nil {
			//					logger.Error(MarshalJSONErr)
			//					return
			//				}
			//				//把jsondata反序列化后进行，自由币判断运算
			//				klineData := huobi.SubData{}
			//				err = json.Unmarshal(jsondata, &klineData)
			//				if err != nil {
			//					logger.Error(err)
			//					return
			//				}
			//				//自由币换算
			//				tranData := logic.TranDecimalScale2(msg, klineData)
			//				//结构体序列化后返回
			//				data, MarshalErr := json.Marshal(tranData)
			//				if MarshalErr != nil {
			//					logger.Error(MarshalErr)
			//					return
			//				}
			//				if len(data) <= 0 {
			//					logger.Info("空数据，不推送:data")
			//					//logger.Info(data)
			//					return
			//				}
			//				//返回数据给用户
			//				wsConn.Mux.Lock()
			//				err = wsConn.Conn.WriteMessage(mt, data)
			//				//logger.Info(data)
			//				wsConn.Mux.Unlock()
			//				//time.Sleep(2*time.Second)
			//				if err != nil {
			//					logger.Error(err)
			//					wsErr := ws.Close()
			//					if wsErr != nil {
			//						logger.Error(wsErr)
			//						return
			//					}
			//
			//				}
			//
			//			})
			//			go func() {
			//				time.Sleep(60 * time.Second)
			//				market.Unsubscribe(newMsg)
			//			}()
			//		}()
			//		market.Loop()
			//
			//	}
			//
			//}()
		}

	}
}
