package echo

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goapi/language"
	"goapi/pkg/logger"
)

// Rjson 成功返回封装 参数 data interface{} 类型为可接受任意类型
func Rjson(c *gin.Context, result interface{}, msg string, code string, success bool) {
	reqId, _ := c.Get("Tracking-Id")
	var rdata map[string]interface{}
	if len(msg) > 0 {
		rdata = gin.H{
			"reqId":   reqId,
			"code":    code,
			"success": success,
			"result":  result,
			"msg":     msg,
		}
	} else {
		rdata = gin.H{
			"reqId":   reqId,
			"code":    code,
			"success": success,
			"result":  result,
			"msg":     "success.",
		}
	}
	c.JSON(200, rdata)
	return
}

// Error  错误返回封装
func Error(c *gin.Context, code string, msg string) {
	logger.Logger.WithOptions(zap.AddCallerSkip(1)).Info("返回错误", zap.Any(code, msg))
	// todo 语言包
	if len(msg) <= 0 {
		code, msg = language.Lang(c.Request.Header.Get("Language")).GetErrorCode(code)
	} else {
		code, _ = language.Lang(c.Request.Header.Get("Language")).GetErrorCode(code)
	}
	logger.Logger.WithOptions(zap.AddCallerSkip(1)).Info("返回错误", zap.Any(code, msg))
	Rjson(c, []interface{}{}, msg, code, false)
}

// Success  错误返回封装
func Success(c *gin.Context, result interface{}, msg string, code string) {
	code = "200"
	Rjson(c, result, msg, code, true)
}
