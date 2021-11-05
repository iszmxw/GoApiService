package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/echo"
	"goapi/pkg/mysql"
	"goapi/pkg/validator"
	"strconv"
)

type TradeController struct {
	BaseController
}

// SysTypeHandler 获取类型
func (h *TradeController) SysTypeHandler(c *gin.Context) {
	var (
		data   []response.GlobalsTypes
		result []map[string]interface{}
	)
	mysql.DB.Debug().Model(models.Globals{}).Where("fields IN ?", []string{"omni_wallet_address", "erc20_wallet_address", "trc20_wallet_address"}).Find(&data)
	if len(data) > 0 {
		for _, val := range data {
			arr := cmap.New().Items()
			switch val.Fields {
			case "omni_wallet_address":
				arr["OMNI"] = val.Value
				break
			case "erc20_wallet_address":
				arr["ERC20"] = val.Value
				break
			case "trc20_wallet_address":
				arr["TRC20"] = val.Value
				break
			}
			result = append(result, arr)
		}
	}
	echo.Success(c, result, "", "")
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
	AddData.TradingPairType = params.TradingPairType
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
		Fees        models.WithdrawalFees
		MinAmount   models.MinAmount
	)
	_ = c.Bind(&params) // 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}

	DB := mysql.DB.Debug()
	DB.Debug().Model(models.Globals{}).Where("fields", "withdrawal_fees").Find(&Fees)
	DB.Debug().Model(&models.Globals{}).Where("fields", "min_amount").Find(&MinAmount)
	if Fees.Id <= 0 || Fees.Value <= 0 {
		echo.Error(c, "withdrawalFeesIsNotExist", "")
		return
	}
	if MinAmount.Id <= 0 || MinAmount.Value <= 0 {
		echo.Error(c, "MinAmountIsNotExist", "")
		return
	}
	userId, _ := c.Get("user_id")
	whereUsersWallet := cmap.New().Items()
	whereUsersWallet["user_id"] = userId
	whereUsersWallet["type"] = params.Type // 1现货 2合约
	whereUsersWallet["trading_pair_id"] = params.TradingPairId
	DB.Model(models.UsersWallet{}).Where(whereUsersWallet).Scan(&UsersWallet)
	if UsersWallet.Id <= 0 {
		echo.Error(c, "UsersWalletIsNotExist", "")
		return
	}
	result := cmap.New().Items()
	result["withdrawal_fees"] = Fees.Value
	result["min_amount"] = MinAmount.Value
	result["balance"] = UsersWallet.Available // 可用余额
	echo.Success(c, result, "", "")

}

// WalletAddressAddHandler 添加提现地址
func (h *TradeController) WalletAddressAddHandler(c *gin.Context) {
	var (
		params  requests.WalletAddressAdd
		AddData models.WalletAddress
	)
	_ = c.Bind(&params) // 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	DB := mysql.DB.Debug().Begin()
	userInfo, _ := c.Get("user")
	userId, _ := c.Get("user_id")
	AddData.UserId = userId.(int)
	AddData.Email = userInfo.(map[string]interface{})["email"].(string)
	AddData.Name = params.Name
	AddData.Pact = params.Pact
	AddData.Address = params.Address
	cErr := DB.Model(models.WalletAddress{}).Create(&AddData).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}
	DB.Commit()
	echo.Success(c, "", "", "")

}

// WalletAddressDelHandler 删除提币地址
func (h *TradeController) WalletAddressDelHandler(c *gin.Context) {
	var (
		params  requests.WalletAddressDel
		DelData models.WalletAddress
	)
	_ = c.Bind(&params) // 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	DB := mysql.DB.Debug().Begin()
	userId, _ := c.Get("user_id")
	cErr := DB.Model(models.WalletAddress{}).
		Where("id", params.Id).
		Where("user_id", userId).
		Delete(&DelData).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}
	DB.Commit()
	echo.Success(c, "", "", "")

}

// WalletAddressListHandler 提币地址列表
func (h *TradeController) WalletAddressListHandler(c *gin.Context) {
	var result []models.WalletAddress
	userId, _ := c.Get("user_id")
	mysql.DB.Debug().Model(models.WalletAddress{}).Where("user_id", userId).Find(&result)
	echo.Success(c, result, "", "")
}

// WithdrawHandler 提币
func (h *TradeController) WithdrawHandler(c *gin.Context) {
	var (
		params        requests.AddWithdraw
		AddData       models.Withdraw
		TradingPair   response.TradingPair
		WalletAddress response.WalletAddress
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
	DB.Model(models.TradingPair{}).Where("id", "1").Find(&TradingPair)
	if TradingPair.Id <= 0 {
		echo.Error(c, "TradingPairIsNotExist", "")
		return
	}
	// 检查现货钱包余额是否充足
	var UsersWallet response.UsersWallet
	DB.Model(models.UsersWallet{}).
		Where("trading_pair_id", "1").
		Where("type", "1").
		Where("user_id", userId).
		Find(&UsersWallet)
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
	DB.Model(models.WalletAddress{}).Where("id", params.AddressId).Find(&WalletAddress)
	// 收集添加数据
	AddData.UserId = userId.(int)
	AddData.Email = userInfo.(map[string]interface{})["email"].(string)
	AddData.Address = WalletAddress.Address
	AddData.TradingPairId = "1"
	AddData.TradingPairName = TradingPair.Name
	AddData.Type = params.Type
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

// GetCurrencyHandler 获取申购币种

func (h *TradeController) GetCurrencyHandler(c *gin.Context) {
	DB := mysql.DB.Debug()
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
	ApplyBuy.GetCurrencyName = ApplyBuySetup.Name  // 申购购买币种名称
	ApplyBuy.TradingPairId = 1                     // 交易对id
	ApplyBuy.TradingPairName = "USDT"              // 交易对名称
	ApplyBuy.IssuePrice = ApplyBuySetup.IssuePrice // 发行价 1 = 多少个USDT
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
