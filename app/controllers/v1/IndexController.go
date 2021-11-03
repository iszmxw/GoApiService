package v1

import (
	"github.com/gin-gonic/gin"
	"goapi/app/models"
	"goapi/app/response"
	"goapi/pkg/echo"
	"goapi/pkg/mysql"
)

type IndexController struct {
	BaseController
}

// SystemInfoHandler 系统信息
func (h *IndexController) SystemInfoHandler(c *gin.Context) {
	var GlobalsTypes1 response.GlobalsTypes
	var GlobalsTypes2 response.GlobalsTypes
	var GlobalsTypes3 response.GlobalsTypes
	mysql.DB.Model(models.Globals{}).Where(map[string]interface{}{"fields": "download_link"}).Find(&GlobalsTypes1)
	mysql.DB.Model(models.Globals{}).Where(map[string]interface{}{"fields": "kf_address"}).Find(&GlobalsTypes2)
	mysql.DB.Model(models.Globals{}).Where(map[string]interface{}{"fields": "kf2_address"}).Find(&GlobalsTypes3)
	result := gin.H{
		"download_link": GlobalsTypes1.Value,
		"kf_address":    GlobalsTypes2.Value,
		"kf2_address":   GlobalsTypes3.Value,
	}
	echo.Success(c, result, "", "")
}

// SysCurrencyHandler 首页-币种列表
func (h *IndexController) SysCurrencyHandler(c *gin.Context) {
	var currencyList []response.Currency
	DB := mysql.DB.Debug()
	DB.Model(models.Currency{}).
		Select(models.Prefix("$_currency.*,$_trading_pair.name as trading_pair_name")).
		Joins(models.Prefix("left join $_trading_pair on $_trading_pair.id=$_currency.trading_pair_id")).
		Find(&currencyList)
	echo.Success(c, currencyList, "", "")
}

// TradingPairHandler 首页-交易对列表
func (h *IndexController) TradingPairHandler(c *gin.Context) {
	var TradingPair []response.TradingPair
	DB := mysql.DB.Debug()
	DB.Model(models.TradingPair{}).Find(&TradingPair)
	echo.Success(c, TradingPair, "", "")
}
