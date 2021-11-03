package v1

import (
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/echo"
	"goapi/pkg/mysql"
	"goapi/pkg/validator"
)

type AssetsStreamController struct {
	BaseController
}

// AssetsStreamHandler 个人资产
func (h *AssetsStreamController) AssetsStreamHandler(c *gin.Context) {
	userId, _ := c.Get("user_id")
	where := cmap.New().Items()
	where["user_id"] = userId
	var result []response.UsersWallet
	DB := mysql.DB.Debug()
	DB.Model(models.UsersWallet{}).Where(where).Order("id DESC").Find(&result)
	echo.Success(c, result, "", "")
}

// AssetsTypeHandler 资产类型，获取单个币种余额
func (h *AssetsStreamController) AssetsTypeHandler(c *gin.Context) {
	// 初始化数据模型结构体
	var params requests.TradingPair
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userId, _ := c.Get("user_id")
	where := cmap.New().Items()
	where["user_id"] = userId
	where["trading_pair_id"] = params.TradingPairId
	var result response.UsersWallet
	DB := mysql.DB.Debug()
	DB.Model(models.UsersWallet{}).Where(where).Find(&result)
	if result.Id <= 0 {
		echo.Error(c, "CurrencyIsExist", "")
		return
	}
	echo.Success(c, result, "", "")
}

// TransferHandler 划转
func (h *AssetsStreamController) TransferHandler(c *gin.Context) {
	userId, _ := c.Get("user_id")
	userInfo, _ := c.Get("user")
	where := cmap.New().Items()
	where["id"] = userId
	echo.Success(c, userInfo, "ok", "")
}
