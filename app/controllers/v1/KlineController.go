package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"goapi/pkg/logger"
	"net/http"
	"strings"
	"sync"
	"time"
)

// k线图服务

type KlineController struct {
	BaseController
}

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsConn 声明并发安全的ws
type WsConn struct {
	*websocket.Conn
	Mux sync.RWMutex
}

// WsHandler socket
func (h *KlineController) WsHandler(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, CloseErr := upGrader.Upgrade(c.Writer, c.Request, nil)
	if CloseErr != nil {
		logger.Error(CloseErr)
	}
	wsConn := &WsConn{
		ws,
		sync.RWMutex{},
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			logger.Error(err)
		}
	}(ws) //返回前关闭
	for {
		market, err := huobiapi.NewMarket()
		if err != nil {
			//logger.Info(111)
			logger.Error(err)
			//logger.Info(666)
		}
		//读取ws中的数据
		mt, message, err := wsConn.Conn.ReadMessage()
		if err != nil {
			//logger.Info(666)
			marketErr := market.Close()
			if marketErr != nil {
				//logger.Info("关闭连接失败1")
				logger.Error(marketErr)
				//logger.Info("关闭连接失败2")
				return
			} else {
				logger.Info("关闭成功")
			}
			logger.Error(err)
			//logger.Info(666)
			break
		}
		//对数据进行切割，读取参数
		//如果请求的是market.ethbtc.kline.5min,订阅这条信息，然后再返回
		msg := string(message)
		newMsg := string([]byte(msg)[1 : len([]byte(msg))-1])
		//打印请求参数
		//logger.Info(newMsg)

		if strings.Contains(msg, "1min") || strings.Contains(msg, "step1") {
			go func() {
				for {
					data, GetDataByKeyErr := logic.GetDataByKey(msg)
					//修改，当拿不到key重新订阅，10秒订阅一次
					if GetDataByKeyErr == redis.Nil {
						logger.Error(errors.New(msg + "：key不存在，准备开始缓存"))
						StartSetKlineDataErr := logic.StartSetKlineData()
						if StartSetKlineDataErr != nil {
							logger.Error(StartSetKlineDataErr)
							return
						}
						time.Sleep(10 * time.Second)
					}
					websocketData := utils.Strval(data)
					if len(websocketData) <= 0 {
						logger.Info("空数据，不推送:websocketData")
						//logger.Info(websocketData)
						return
					}
					wsConn.Mux.Lock()
					err = wsConn.Conn.WriteMessage(mt, []byte(websocketData))
					//logger.Info(websocketData)
					wsConn.Mux.Unlock()
					if err != nil {
						logger.Error(err)
						wsErr := ws.Close()
						if wsErr != nil {
							logger.Error(wsErr)
							return
						}
						return
					}
					time.Sleep(time.Second * 2)
				}

			}()
		} else {
			//写入ws数据
			go func() {
				for {

					go func() {
						err = market.Subscribe(newMsg, func(topic string, hjson *huobiapi.JSON) {
							//logger.Info(msg)
							if err != nil {
								logger.Error(err)
							}
							//订阅成功
							//logger.Info("订阅成功")
							//120后自动取消订阅
							go func() {
								time.Sleep(60 * time.Minute)
								//logger.Info("取消订阅成功")
								market.Unsubscribe(newMsg)
								//market.ReceiveTimeout

							}()

							// 收到数据更新时回调
							//logger.Info(topic)
							//logger.Info(hjson)
							jsondata, MarshalJSONErr := hjson.MarshalJSON()
							if err != nil {
								logger.Error(MarshalJSONErr)
								return
							}
							//把jsondata反序列化后进行，自由币判断运算
							klineData := huobi.SubData{}
							err = json.Unmarshal(jsondata, &klineData)
							if err != nil {
								logger.Error(err)
								return
							}
							//自由币换算
							tranData := logic.TranDecimalScale2(msg, klineData)
							//结构体序列化后返回
							data, MarshalErr := json.Marshal(tranData)
							if MarshalErr != nil {
								logger.Error(MarshalErr)
								return
							}
							if len(data) <= 0 {
								logger.Info("空数据，不推送:data")
								//logger.Info(data)
								return
							}
							//返回数据给用户
							wsConn.Mux.Lock()
							err = wsConn.Conn.WriteMessage(mt, data)
							//logger.Info(data)
							wsConn.Mux.Unlock()
							//time.Sleep(2*time.Second)
							if err != nil {
								logger.Error(err)
								wsErr := ws.Close()
								if wsErr != nil {
									logger.Error(wsErr)
									return
								}

							}

						})
						go func() {
							time.Sleep(60 * time.Second)
							market.Unsubscribe(newMsg)
						}()
					}()
					market.Loop()

				}

			}()
		}

	}
}
