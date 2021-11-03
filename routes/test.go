package routes

import (
	"github.com/gin-gonic/gin"
	"goapi/app/controllers/test"
)

// RegisterTestRoutes 注册测试路由
func RegisterTestRoutes(router *gin.RouterGroup) {
	// 路由分组 客户端 模块
	ApiRoute := router.Group("/api")
	{
		Test := new(test.Controller)
		// 测试 redis
		ApiRoute.Any("/redis", Test.SetHandler)
		ApiRoute.Any("/dy", Test.DyHandler)
	}
}
