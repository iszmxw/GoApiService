package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goapi/pkg/echo"
	"goapi/pkg/logger"
	"goapi/pkg/redis"
	"goapi/pkg/request"
	"strconv"
)

// Client 定义客户端中间件
func Client() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置语言，默认英文
		SetLang(c, "en")
		logger.WithContext(c).Info("请求信息", zap.Skip())
		path := c.FullPath()
		switch path {
		case "/v1/api/user/login":
		case "/v1/api/user/send_email_register":
		case "/v1/api/user/verify_register":
		case "/v1/api/user/send_email_retrieve":
		case "/v1/api/user/reset_verify":
		case "/v1/api/user/reset_password":
		case "/v1/api/index/trading_pair":
		case "/v1/api/index/banner": // 首页轮播图
		case "/v1/api/index/sys_currency":
		case "/v1/api/index/system_info":
			// 继续往下面执行
			c.Next()
			break
		default:
			CheckLogin(c)
			c.Next()
			break
		}
	}
}

// CheckLogin 检测登录
func CheckLogin(c *gin.Context) {
	var (
		strInfo string
		err     error
	)
	// 获取 "token"
	tokenString := request.GetParam(c, "token")
	if len(tokenString) <= 0 || tokenString == "<nil>" {
		logger.Error(errors.New("未检测到token"))
		echo.Error(c, "LoginError", "")
		c.Abort()
		return
	}
	strInfo, err = redis.Get(tokenString)
	logger.Info(strInfo)
	if err != nil {
		echo.Error(c, "LoginError", "")
		c.Abort()
		return
	}
	// json字符串数组,转换成切片
	var user map[string]interface{}
	multiErr := json.Unmarshal([]byte(strInfo), &user)
	if multiErr != nil {
		logger.Error(errors.New("转换出错"))
		logger.Error(multiErr)
		return
	}
	// id 转换为整型
	userId, strErr := strconv.Atoi(fmt.Sprintf("%.0f", user["id"].(float64)))
	if strErr != nil {
		logger.Error(err)
	}
	user["id"] = userId
	// 保存用户到 上下文
	c.Set("user", user)
	c.Set("language", user["language"])
	c.Set("user_id", userId)
	// 再次确认设置语言
	SetLang(c, "zh")
	logger.Info("=================middlewares=====================")
	// 继续往下面执行
	c.Next()
}

// SetLang 设置语言
func SetLang(c *gin.Context, lang string) {
	language := c.Request.Header.Get("Language")
	if len(language) <= 0 {
		// 默认中文
		c.Request.Header.Set("Language", lang)
	}
	l, _ := c.Get("language")
	switch l {
	case 1:
		// 中文
		c.Request.Header.Set("Language", "zh")
		break
	case 2:
		// 英文
		c.Request.Header.Set("Language", "en")
		break
	case 3:
		// 日语
		c.Request.Header.Set("Language", "jp")
		break
	}
}
