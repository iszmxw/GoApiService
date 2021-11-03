package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goapi/pkg/echo"
	"goapi/pkg/redis"
)

type Controller struct {
}

// SetHandler 登录
func (h *Controller) SetHandler(c *gin.Context) {

	r, err := redis.Add("123456", "test", 10) // 缓存两个小时过期
	if err != nil {
		fmt.Println(err.Error())
	}
	// 这里sleep是为了防止main方法直接推出
	echo.Success(c, r, "", "")
}

// DyHandler 订阅消息
func (h *Controller) DyHandler(c *gin.Context) {
	//sub := redis.SubExpireEvent()
	//for {
	//	msg := <-sub.Channel()
	//	fmt.Println("Channel ", msg.Channel)
	//	fmt.Println("pattern ", msg.Pattern)
	//	fmt.Println("pattern ", msg.Payload)
	//}
}
