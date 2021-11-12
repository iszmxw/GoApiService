package v1

import (
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/echo"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/validator"
	"time"
)

type OrderController struct {
	BaseController
}

// TypeHandler 流水类型
func (h *OrderController) TypeHandler(c *gin.Context) {
	type Rdata struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var arr []Rdata
	// 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	arr = append(arr, Rdata{
		Id:   1,
		Name: "充值",
	})
	arr = append(arr, Rdata{
		Id:   2,
		Name: "提现",
	})
	arr = append(arr, Rdata{
		Id:   3,
		Name: "划转",
	})
	arr = append(arr, Rdata{
		Id:   4,
		Name: "快捷买币",
	})
	arr = append(arr, Rdata{
		Id:   5,
		Name: "空投",
	})
	arr = append(arr, Rdata{
		Id:   6,
		Name: "现货",
	})
	arr = append(arr, Rdata{
		Id:   7,
		Name: "合约",
	})
	arr = append(arr, Rdata{
		Id:   8,
		Name: "期权",
	})
	arr = append(arr, Rdata{
		Id:   9,
		Name: "手续费",
	})
	echo.Success(c, arr, "", "")
}

// ListHandler 流水列表
func (h *OrderController) ListHandler(c *gin.Context) {
	var (
		params   requests.ListAssetsStream
		pageList models.PageList // 返回数据
		result   []response.WalletStream
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
	DB := mysql.DB.Debug()
	DB = DB.Model(models.WalletStream{})
	// 时间段搜索
	if len(params.StartTime) > 0 && len(params.EndTime) > 0 {
		//格式化字符串为时间
		StartTime, StartTimeErr := time.Parse("2006-01-02", params.StartTime)
		EndTime, EndTimeErr := time.Parse("2006-01-02", params.EndTime)
		if StartTimeErr != nil || EndTimeErr != nil {
			echo.Error(c, "SearchTimeErr", "")
			return
		}
		DB = DB.Where("created_at >= ?", StartTime).Where("created_at <= ?", EndTime)
	}
	table := DB.Where(where)
	table.Count(&pageList.Total)
	// 设置分页参数
	pageList.CurrentPage = int64(params.Page)
	pageList.PageSize = int64(params.Limit)
	models.InitPageList(&pageList)
	// order by
	table = table.Order("id desc").Offset(int(pageList.Offset)).Limit(int(pageList.PageSize))
	if err := table.Scan(&result).Error; err != nil {
		// 记录错误
		logger.Error(err)
	} else {
		pageList.Data = result
	}
	echo.Success(c, pageList, "ok", "")
}
