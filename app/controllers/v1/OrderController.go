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
	"time"
)

type OrderController struct {
	BaseController
}

// TypeHandler 流水类型
func (h *OrderController) TypeHandler(c *gin.Context) {
	rdata := cmap.New().Items()
	// 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	rdata["1"] = "充值"
	rdata["2"] = "提现"
	rdata["3"] = "划转"
	rdata["4"] = "快捷买币"
	rdata["5"] = "空投"
	rdata["6"] = "现货"
	rdata["7"] = "合约"
	rdata["8"] = "期权"
	rdata["9"] = "手续费"
	echo.Success(c, rdata, "", "")
}

// ListHandler 流水列表
func (h *OrderController) ListHandler(c *gin.Context) {
	var (
		params requests.ListAssetsStream
		list   []response.WalletStream
	)
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userId, _ := c.Get("user_id")
	//userInfo, _ := c.Get("user")
	where := cmap.New().Items()
	where["user_id"] = userId
	if len(params.TradingPairId) > 0 {
		where["trading_pair_id"] = params.TradingPairId
	}
	if len(params.OrderType) > 0 {
		where["type"] = params.OrderType
	}
	//获取两天前的时间
	currentTime := time.Now()
	DB := mysql.DB.Debug()
	DB = DB.Model(models.WalletStream{})
	if len(params.Time) > 0 {
		var oldTime string
		switch params.Time {
		case "7":
			oldTime = currentTime.AddDate(0, 0, -7).Format("2006-01-02 15:04:05") // 前七天时间
			//oldTime 的结果为go的时间time类型，2018-09-25 13:24:58.287714118 +0000 UTC
			break
		case "15":
			oldTime = currentTime.AddDate(0, 0, -15).Format("2006-01-02 15:04:05") // 前15天时间
			//oldTime 的结果为go的时间time类型，2018-09-25 13:24:58.287714118 +0000 UTC
			break
		case "30":
			oldTime = currentTime.AddDate(0, 0, -30).Format("2006-01-02 15:04:05") // 前30天时间
			break
		}
		if len(oldTime) > 0 {
			DB.Where("created_at > ?", oldTime)
		}
	}
	DB.Where(where).Order("id desc").Find(&list)
	echo.Success(c, list, "ok", "")
}
