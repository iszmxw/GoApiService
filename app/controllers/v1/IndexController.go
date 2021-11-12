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

type IndexController struct {
	BaseController
}

// SystemInfoHandler 系统信息
func (h *IndexController) SystemInfoHandler(c *gin.Context) {
	var (
		GlobalsTypes1 response.GlobalsTypes
		GlobalsTypes2 response.GlobalsTypes
		GlobalsTypes3 response.GlobalsTypes
	)
	mysql.DB.Model(models.Globals{}).Where("fields", "download_link").Find(&GlobalsTypes1)
	mysql.DB.Model(models.Globals{}).Where("fields", "kfh5_address").Find(&GlobalsTypes2)
	mysql.DB.Model(models.Globals{}).Where("fields", "kfapp_address").Find(&GlobalsTypes3)
	result := gin.H{
		"download_link": GlobalsTypes1.Value, // APP下载地址
		"kfh5_address":  GlobalsTypes2.Value, // 客服H5地址
		"kfapp_address": GlobalsTypes3.Value, // 客服APP地址
	}
	echo.Success(c, result, "", "")
}

// BannerHandler 轮播图
func (h *IndexController) BannerHandler(c *gin.Context) {
	var (
		params requests.BannerType
		Banner []response.Banner
	)
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	where := cmap.New().Items()
	if len(params.Type) > 0 {
		where["type"] = params.Type
	}
	mysql.DB.Model(models.Banner{}).Where(where).Find(&Banner)
	echo.Success(c, Banner, "", "")
}

// SysCurrencyHandler 首页-币种列表
func (h *IndexController) SysCurrencyHandler(c *gin.Context) {
	var (
		currencyList []response.Currency
		params       requests.SysCurrencyType
	)
	_ = c.Bind(&params)
	DB := mysql.DB.Debug()
	Query := DB.Model(models.Currency{}).
		Where("is_hidden", "1").
		Select(models.Prefix("$_currency.*,$_trading_pair.name as trading_pair_name")).
		Joins(models.Prefix("left join $_trading_pair on $_trading_pair.id=$_currency.trading_pair_id"))
	// 搜索条件
	if len(params.Type) > 0 {
		Query = Query.Where("FIND_IN_SET(?,type)", params.Type)
	}
	Query.Find(&currencyList)
	echo.Success(c, currencyList, "", "")
}

// TradingPairHandler 首页-交易对列表
func (h *IndexController) TradingPairHandler(c *gin.Context) {
	var TradingPair []response.TradingPair
	DB := mysql.DB.Debug()
	DB.Model(models.TradingPair{}).Find(&TradingPair)
	echo.Success(c, TradingPair, "", "")
}
