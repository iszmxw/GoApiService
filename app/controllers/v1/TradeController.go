package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/echo"
	"goapi/pkg/helpers"
	"goapi/pkg/mysql"
	"goapi/pkg/validator"
	"strconv"
)

type TradeController struct {
	BaseController
}

// SsyTypeHandler 获取类型
func (h *TradeController) SsyTypeHandler(c *gin.Context) {
	var data []response.GlobalsTypes
	mysql.DB.Debug().Model(models.Globals{}).Where("fields IN ?", []string{"omni_wallet_address", "erc20_wallet_address", "trc20_wallet_address"}).Find(&data)
	echo.Success(c, data, "", "")
}

// ReChargeHandler 充值
func (h *TradeController) ReChargeHandler(c *gin.Context) {
	var (
		params  requests.Recharge
		AddData models.Recharge
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
	userInfo, _ := c.Get("user")
	DB := mysql.DB.Begin()
	// 收集添加数据
	AddData.UserId = userId.(int)
	AddData.Email = userInfo.(map[string]interface{})["email"].(string)
	AddData.Address = params.Address
	AddData.Type = params.Type
	AddData.TradingPairId = params.TradingPairId
	AddData.TradingPairName = params.TradingPairName
	AddData.RechargeNum = params.RechargeNum
	AddData.Status = "0"
	// 收集添加数据
	cErr := DB.Debug().Model(models.Recharge{}).Create(&AddData).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}
	DB.Commit()
	echo.Success(c, AddData, "ok", "")
}

// ReChargeLogHandler 充值记录
func (h *TradeController) ReChargeLogHandler(c *gin.Context) {
	var (
		request requests.ListRecharge // 接收参数
		lists   models.PageList       // 返回数据
		DB      models.Recharge       // 数据模型
	)
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&request)
	// 数据验证
	vErr := validator.Validate.Struct(request)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, request, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	email := userInfo.(map[string]interface{})["email"].(string)
	where := cmap.New().Items()
	where["email"] = email
	if len(request.TradingPairName) > 0 {
		where["trading_pair_name"] = request.TradingPairName
	}
	if len(request.Status) > 0 {
		where["status"] = request.Status
	}
	// 绑定接收的 json 数据到结构体中
	DB.GetPaginate(where, request.OrderBy, int64(request.Page), int64(request.Limit), &lists)
	echo.Success(c, lists, "ok", "")
}

// GetWithdrawConfigHandler 获取提现配置信息
func (h *TradeController) GetWithdrawConfigHandler(c *gin.Context) {
	var (
		params      requests.TradingPair
		UsersWallet response.UsersWallet
	)
	_ = c.Bind(&params) // 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	withdrawalFees := map[string]interface{}{}
	DB := mysql.DB.Debug()
	DB.Debug().Model(&models.Globals{}).Where("fields", "withdrawal_fees").Scan(withdrawalFees)
	if len(withdrawalFees) <= 0 || withdrawalFees["id"].(uint32) <= 0 || len(withdrawalFees["value"].(string)) <= 0 {
		echo.Error(c, "withdrawalFeesIsNotExist", "")
		return
	}
	userId, _ := c.Get("user_id")
	whereUsersWallet := cmap.New().Items()
	whereUsersWallet["user_id"] = userId
	whereUsersWallet["trading_pair_id"] = params.TradingPairId
	DB.Model(models.UsersWallet{}).Where(whereUsersWallet).Scan(&UsersWallet)
	if UsersWallet.Id <= 0 {
		echo.Error(c, "UsersWalletIsNotExist", "")
		return
	}
	result := cmap.New().Items()
	result["withdrawal_fees"] = withdrawalFees["value"]
	result["balance"] = UsersWallet.Available // 可用余额
	echo.Success(c, result, "", "")

}

// WithdrawHandler 提币
func (h *TradeController) WithdrawHandler(c *gin.Context) {
	var (
		params      requests.AddWithdraw
		AddData     models.Withdraw
		TradingPair response.TradingPair
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
	userInfo, _ := c.Get("user")
	DB := mysql.DB.Debug().Begin()
	DB.Model(models.TradingPair{}).Where("id", params.TradingPairId).Find(&TradingPair)
	if TradingPair.Id <= 0 {
		echo.Error(c, "TradingPairIsNotExist", "")
		return
	}
	// 检查钱包余额是否充足
	var UsersWallet response.UsersWallet
	DB.Model(models.UsersWallet{}).Where(map[string]interface{}{"user_id": userId, "trading_pair_id": params.TradingPairId}).Find(&UsersWallet)
	WithdrawNum, _ := strconv.ParseFloat(params.WithdrawNum, 64)
	if WithdrawNum > UsersWallet.Available {
		echo.Error(c, "InsufficientBalance", "")
		return
	}
	var WithdrawalFees response.WithdrawalFees
	DB.Model(models.Globals{}).Where("fields", "withdrawal_fees").Find(&WithdrawalFees)
	if WithdrawalFees.Value < 0 {
		echo.Error(c, "WithdrawalFeesIsError", "")
		return
	}

	// 收集添加数据
	AddData.UserId = userId.(int)
	AddData.Email = userInfo.(map[string]interface{})["email"].(string)
	AddData.Address = params.Address
	AddData.TradingPairId = params.TradingPairId
	AddData.TradingPairName = TradingPair.Name
	AddData.Type = helpers.StringToInt(params.Type) + 1
	AddData.WithdrawNum = WithdrawNum
	AddData.HandlingFee = WithdrawalFees.Value
	AddData.ActuallyArrived = WithdrawNum - (WithdrawNum * (WithdrawalFees.Value / 100))
	AddData.Status = "0"
	cErr := DB.Debug().Model(AddData).Create(&AddData).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}
	DB.Commit()
	echo.Success(c, AddData, "ok", "")
}

// WithdrawLogHandler 提币记录
func (h *TradeController) WithdrawLogHandler(c *gin.Context) {
	var (
		request requests.ListWithdraw // 接收参数
		lists   models.PageList       // 返回数据
		DB      models.Withdraw       // 数据模型
	)
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&request)
	// 数据验证
	vErr := validator.Validate.Struct(request)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, request, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	email := userInfo.(map[string]interface{})["email"].(string)
	where := cmap.New().Items()
	where["email"] = email
	if len(request.TradingPairName) > 0 {
		where["trading_pair_name"] = request.TradingPairName
	}
	if len(request.Status) > 0 {
		where["status"] = request.Status
	}
	// 绑定接收的 json 数据到结构体中
	DB.GetPaginate(where, request.OrderBy, int64(request.Page), int64(request.Limit), &lists)
	echo.Success(c, lists, "ok", "")
}

// GetCurrencyHandler  todo:over
func (h *TradeController) GetCurrencyHandler(c *gin.Context) {
	userId, _ := c.Get("user_id")
	where := cmap.New().Items()
	where["id"] = userId
	DB := mysql.DB.Debug().Begin()
	list := new([]response.ApplyBuySetup)
	DB.Model(models.ApplyBuySetup{}).Find(&list)
	echo.Success(c, list, "", "")
}

// SubmitApplyBuyHandler 交易相关-申购
func (h *TradeController) SubmitApplyBuyHandler(c *gin.Context) {
	userInfo, _ := c.Get("user")
	var params requests.SubmitApplyBuy
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	user := userInfo.(map[string]interface{})
	ApplyBuy := models.ApplyBuy{}
	ApplyBuy.Email = user["email"].(string)         // 邮件
	ApplyBuy.GetCurrencyId = params.GetCurrencyId   // 申购购买币种id
	ApplyBuy.GetCurrencyNum = params.GetCurrencyNum // 申购购买数量
	DB := mysql.DB.Debug().Begin()
	var ApplyBuySetup models.ApplyBuySetup
	DB.Model(ApplyBuySetup).Where("id", params.GetCurrencyId).Find(&ApplyBuySetup)
	if ApplyBuySetup.Id <= 0 {
		echo.Error(c, "ApplyBuySetupIsNotExist", "")
		return
	}
	ApplyBuy.GetCurrencyName = ApplyBuySetup.Name            // 申购购买币种名称
	ApplyBuy.TradingPairId = ApplyBuySetup.TradingPairId     // 交易对id
	ApplyBuy.TradingPairName = ApplyBuySetup.TradingPairName // 交易对名称
	ApplyBuy.Ratio = ApplyBuySetup.Ratio                     // 购买比例
	dbErr := DB.Model(ApplyBuy).Create(&ApplyBuy).Error
	if dbErr != nil {
		DB.Rollback()
		fmt.Printf("添加数据失败%v\n", dbErr.Error())
		echo.Error(c, "AddError", dbErr.Error())
		return
	}
	DB.Commit()

	echo.Success(c, ApplyBuy, "", "")
}
