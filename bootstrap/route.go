package bootstrap

import (
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"goapi/app/middlewares/v1"
	_ "goapi/docs"
	conf "goapi/pkg/config"
	"goapi/routes"
)

// SetupRoute 路由初始化
func SetupRoute() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// 支持wss
	if conf.GetString("app.https") == "1" {
		router.Use(v1.TlsHandler())
	}
	router.Use(v1.TraceLogger()) // 日志追踪
	router.Use(v1.Cors())        // 跨域
	// swagger docs 文档
	router.GET("/docs/*any", gs.WrapHandler(swaggerFiles.Handler))
	// v1 版本
	apiV1 := router.Group("/v1")
	Test := router.Group("/test")
	router.GET("/", func(context *gin.Context) {
		requestId, _ := context.Get("Tracking-Id")
		context.String(200, "Hello World!："+requestId.(string)+"\n\n\n")
		//context.String(200, "下面是所有接口服务：\n\n\n")
		//routers := router.Routes()
		//for _, v := range routers {
		//	context.String(200, fmt.Sprintf("Method：\t%v  URL：\t%v  \t\t\tHandler: \t%v \n", v.Method, v.Path, v.Handler))
		//}
	})
	routes.RegisterWebRoutes(apiV1)
	routes.RegisterTestRoutes(Test)
	return router
}
