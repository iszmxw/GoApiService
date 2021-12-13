package v1

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/skip2/go-qrcode"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/config"
	"goapi/pkg/echo"
	"goapi/pkg/helpers"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/validator"
	"io/ioutil"
	"net/http"
	"strconv"
)

type TradeController struct {
	BaseController
}

// SysTypeHandler 获取类型
func (h *TradeController) SysTypeHandler(c *gin.Context) {
	var (
		data   []response.GlobalsTypes  // 查询系统配置
		result []map[string]interface{} // 响应结果
		Img    string                   // 钱包地址图片
	)
	mysql.DB.Debug().Model(models.Globals{}).Where("fields IN ?", []string{"omni_wallet_address", "erc20_wallet_address", "trc20_wallet_address"}).Find(&data)
	if len(data) > 0 {
		for _, val := range data {
			arr := cmap.New().Items()
			switch val.Fields {
			case "omni_wallet_address":
				arr["name"] = "OMNI"
				arr["address"] = val.Value
				OmniImg, OmniImgerr := qrcode.Encode(val.Value, qrcode.Medium, 256)
				if OmniImgerr != nil {
					Img = ""
				} else {
					Img = "data:image/png;base64," + base64.StdEncoding.EncodeToString(OmniImg)
				}
				arr["resource"] = Img
				break
			case "erc20_wallet_address":
				arr["name"] = "ERC20"
				arr["address"] = val.Value
				ecr20Img, ecr20Imgerr := qrcode.Encode(val.Value, qrcode.Medium, 256)
				if ecr20Imgerr != nil {
					Img = ""
				} else {
					Img = "data:image/png;base64," + base64.StdEncoding.EncodeToString(ecr20Img)
				}
				arr["resource"] = Img
				break
			case "trc20_wallet_address":
				arr["name"] = "TRC20"
				arr["address"] = val.Value
				tcr20Img, tcr20Imgerr := qrcode.Encode(val.Value, qrcode.Medium, 256)
				if tcr20Imgerr != nil {
					Img = ""
				} else {
					Img = "data:image/png;base64," + base64.StdEncoding.EncodeToString(tcr20Img)
				}
				arr["resource"] = Img
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
		result  response.RespData
	)
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	// 银行卡充值
	if params.TopUpType == "2" {
		// 银行账号不能为空
		if len(params.AccountNo) <= 0 {
			echo.Error(c, "ParamAccountNo", "")
			return
		}
		// 银行编码不能为空
		if len(params.BankCode) <= 0 {
			echo.Error(c, "ParamBankCode", "")
			return
		}
		// 检测产品参数，默认为 ThaiP2P
		if len(params.Product) <= 0 {
			params.Product = "ThaiP2P"
		}
		token := c.Request.Header.Get("token")
		url := fmt.Sprintf(config.GetString("app.php_url")+"/api/client/pay/create?amount=%v&account_no=%v&bank_code=%v&product=%v&token=%v", params.RechargeNum, params.AccountNo, params.BankCode, params.Product, token)
		logger.Info(url)
		resp, getErr := http.Get(url)
		body, _ := ioutil.ReadAll(resp.Body)
		if getErr != nil {
			logger.Error(getErr)
			echo.Error(c, "ParamBankErr", "")
			return
		}
		UnmarshalErr := json.Unmarshal(body, &result)
		if UnmarshalErr != nil {
			logger.Error(UnmarshalErr)
			echo.Error(c, "ParsingError", "")
			return
		}
		if result.Data.Success == false {
			logger.Error(errors.New(result.Data.Msg))
			echo.Error(c, "ParamBankErr", "")
			return
		}
		logger.Info(result)
		// 添加收集数据
		AddData.PayId = int(result.Data.Result.(map[string]interface{})["id"].(float64))
	} else {
		params.TopUpType = "1"
	}
	userId, _ := c.Get("user_id")
	userInfo, _ := c.Get("user")
	DB := mysql.DB.Begin()
	// 收集添加数据
	AddData.TopUpType = helpers.StringToInt(params.TopUpType)
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
	// 银行卡充值返回
	if params.TopUpType == "2" {
		echo.Success(c, result.Data.Result, "ok", "")
		return
	} else {
		// USDT 返回
		echo.Success(c, AddData, "ok", "")
		return
	}
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
	if Fees.Id <= 0 {
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
	if len(result) > 0 {
		for k, v := range result {
			switch v.Pact {
			// 协议： 1-OMNI 2-ERC20 3-TRC20
			case "1":
				result[k].Pact = "OMNI"
				break
			case "2":
				result[k].Pact = "ERC20"
				break
			case "3":
				result[k].Pact = "TRC20"
				break
			}
		}
	}
	echo.Success(c, result, "", "")
}

// WithdrawHandler 提币
func (h *TradeController) WithdrawHandler(c *gin.Context) {
	var (
		params         requests.AddWithdraw    // 接收请求参数
		TradingPair    response.TradingPair    // 查询交易对
		WalletAddress  response.WalletAddress  // 查询钱包地址
		UserStatus     response.User           // 查询用户状态
		UsersWallet    response.UsersWallet    // 查询现货钱包
		WithdrawalFees response.WithdrawalFees // 查询提现费率
		AddData        models.Withdraw         // 添加提现数据
		WithdrawStream models.WalletStream     // 提现流水
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
	DB := mysql.DB.Debug()
	DB.Model(models.User{}).Where("id", userId).Find(&UserStatus)
	if UserStatus.Status == "1" {
		echo.Error(c, "UserIsLock", "")
		return
	}
	DB.Model(models.TradingPair{}).Where("id", "1").Find(&TradingPair)
	if TradingPair.Id <= 0 {
		echo.Error(c, "TradingPairIsNotExist", "")
		return
	}
	// 检查现货钱包余额是否充足
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
	DB.Model(models.Globals{}).Where("fields", "withdrawal_fees").Find(&WithdrawalFees)
	if WithdrawalFees.Value < 0 {
		echo.Error(c, "WithdrawalFeesIsError", "")
		return
	}
	DB = DB.Begin()
	// 扣减钱包余额,
	Available := UsersWallet.Available - WithdrawNum
	AvailableErr := DB.Model(models.UsersWallet{}).
		Where("trading_pair_id", "1").
		Where("type", "1").
		Where("user_id", userId).
		Update("available", Available).Error
	if AvailableErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", AvailableErr.Error())
		return
	}

	DB.Model(models.WalletAddress{}).Where("id", params.AddressId).Find(&WalletAddress)
	if len(WalletAddress.Address) <= 0 {
		DB.Rollback()
		echo.Error(c, "WalletAddressErr", "")
		return
	}
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
		logger.Error(cErr)
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}

	// 创建钱包流水
	WithdrawStream.Way = "2"                                                   // 流转方式 1 收入 2 支出
	WithdrawStream.Type = "7"                                                  // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	WithdrawStream.TypeDetail = "5"                                            // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	WithdrawStream.UserId = userId.(int)                                       // 用户id
	WithdrawStream.Email = userInfo.(map[string]interface{})["email"].(string) // 用户邮箱
	WithdrawStream.TradingPairId = 1                                           // 交易对id
	WithdrawStream.TradingPairName = TradingPair.Name                          // 交易对名称
	WithdrawStream.Amount = fmt.Sprintf("%v", WithdrawNum)                     // 流转金额
	WithdrawStream.AmountBefore = UsersWallet.Available                        // 流转前的余额
	WithdrawStream.AmountAfter = Available                                     // 流转后的余额
	WithdrawStreamcErr := DB.Model(models.WalletStream{}).Create(&WithdrawStream).Error
	if WithdrawStreamcErr != nil {
		logger.Error(errors.New(fmt.Sprintf("创建钱包流水失败，%v", WithdrawStreamcErr.Error())))
		DB.Rollback()
		echo.Error(c, "AddError", "")
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
	var (
		params         requests.SubmitApplyBuy // 接收请求参数
		ApplyBuySetup  response.ApplyBuySetup  // 查询申购币种的信息
		GetUsersWallet response.UsersWallet    // 查询申购对应币种的钱包
		UsersWallet    response.UsersWallet    // 查询消费的钱包信息
		UserStatus     response.User           // 查询用户状态
		data           models.WalletStream     // 交易流水创建
		data1          models.WalletStream     // 交易流水创建
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
	user := userInfo.(map[string]interface{})
	ApplyBuy := models.ApplyBuy{}
	ApplyBuy.UserId = helpers.IntToString(userId.(int)) // 用户ID
	ApplyBuy.Email = user["email"].(string)             // 用户邮件
	ApplyBuy.GetCurrencyId = params.GetCurrencyId       // 申购购买币种id
	ApplyBuy.GetCurrencyNum = params.GetCurrencyNum     // 申购购买数量
	DB := mysql.DB.Debug().Begin()
	DB.Model(models.User{}).Where("id", userId).Find(&UserStatus)
	if UserStatus.Status == "1" {
		DB.Rollback()
		echo.Error(c, "UserIsLock", "")
		return
	}
	DB.Model(models.ApplyBuySetup{}).Where("id", params.GetCurrencyId).Find(&ApplyBuySetup)
	if ApplyBuySetup.Id <= 0 {
		echo.Error(c, "ApplyBuySetupIsNotExist", "")
		return
	}
	if ApplyBuySetup.Status != 1 { // 币种状态  0 关闭 1 开启
		echo.Error(c, "ApplyBuySetupStatusIsNotExist", "")
		return
	}
	//// 开启了申购码
	//if ApplyBuySetup.CodeStatus == 1 && params.Code != ApplyBuySetup.Code{ // 申购码开关  0 关闭 1 开启
	//	echo.Error(c, "ApplyBuySetupStatusIsNotExist", "")
	//	return
	//}

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

	// 消费掉的 USDT
	Amount := ApplyBuySetup.IssuePrice * params.GetCurrencyNum

	// 查询用户是否有该钱包
	DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("trading_pair_name", ApplyBuy.GetCurrencyName).
		Where("type", "2").
		Find(&GetUsersWallet)
	if GetUsersWallet.Id <= 0 {
		DB.Rollback()
		echo.Error(c, "UsersWalletIsNotExist", "")
		return
	}
	// 查询用户钱包信息
	DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("trading_pair_id", "1").
		Where("type", "2"). // 钱包类型：1现货 2合约
		Find(&UsersWallet)
	// 用户可用余额不足
	if UsersWallet.Available < 0 || UsersWallet.Id <= 0 {
		logger.Error(errors.New(fmt.Sprintf("UsersWallet.Available: %v < 0 || UsersWallet.Id: %v <= 0", UsersWallet.Available, UsersWallet.Id)))
		DB.Rollback()
		echo.Error(c, "InsufficientBalance", "")
		return
	}

	// 修改钱包余额 （空投消费）
	UpdateUsersWallet := cmap.New().Items()
	UpdateUsersWallet["available"] = UsersWallet.Available - Amount // 修改余额
	if UsersWallet.Available < Amount {
		echo.Error(c, "InsufficientBalance", "")
		return
	}
	editError := DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("type", "2"). // 钱包类型：1现货 2合约
		Where("trading_pair_id", "1").
		Updates(UpdateUsersWallet).Error
	if editError != nil {
		logger.Error(errors.New(fmt.Sprintf("修改钱包余额失败，%v", editError.Error())))
		DB.Rollback()
		echo.Error(c, "AddError", editError.Error())
		return
	}
	// 修改钱包余额

	// 修改钱包余额 （空投收入）
	UpdateUsersWallet1 := cmap.New().Items()
	UpdateUsersWallet1["available"] = GetUsersWallet.Available + Amount // 申购所得的币种
	editError1 := DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("type", "2"). // 钱包类型：1现货 2合约
		Where("trading_pair_name", ApplyBuy.GetCurrencyName).
		Updates(UpdateUsersWallet1).Error
	if editError1 != nil {
		logger.Error(errors.New(fmt.Sprintf("修改钱包余额失败，%v", editError.Error())))
		DB.Rollback()
		echo.Error(c, "AddError", editError1.Error())
		return
	}
	// 修改钱包余额

	// 创建钱包消费流水
	data.Way = "2"                                    // 流转方式 1 收入 2 支出
	data.Type = "5"                                   // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	data.TypeDetail = "6"                             // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	data.UserId = userId.(int)                        // 用户id
	data.Email = user["email"].(string)               // 用户邮箱
	data.TradingPairId = 1                            // 交易对id
	data.TradingPairName = "USDT"                     // 交易对名称
	data.Amount = fmt.Sprintf("%v", Amount)           // 流转金额
	data.AmountBefore = UsersWallet.Available         // 流转前的余额
	data.AmountAfter = UsersWallet.Available - Amount // 流转后的余额
	cErr := DB.Model(models.WalletStream{}).Create(&data).Error
	if cErr != nil {
		logger.Error(errors.New(fmt.Sprintf("创建钱包流水失败，%v", cErr.Error())))
		DB.Rollback()
		echo.Error(c, "AddError", "")
		return
	}

	// 创建钱包收入流水
	data1.Way = "1"                                        // 流转方式 1 收入 2 支出
	data1.Type = "5"                                       // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	data1.TypeDetail = "7"                                 // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	data1.UserId = userId.(int)                            // 用户id
	data1.Email = user["email"].(string)                   // 用户邮箱
	data1.TradingPairId = GetUsersWallet.TradingPairId     // 交易对id
	data1.TradingPairName = GetUsersWallet.TradingPairName // 交易对名称
	data1.Amount = fmt.Sprintf("%v", Amount)               // 流转金额
	data1.AmountBefore = GetUsersWallet.Available          // 流转前的余额
	data1.AmountAfter = GetUsersWallet.Available + Amount  // 流转后的余额
	cErr1 := DB.Model(models.WalletStream{}).Create(&data1).Error
	if cErr1 != nil {
		logger.Error(errors.New(fmt.Sprintf("创建钱包流水失败，%v", cErr1.Error())))
		DB.Rollback()
		echo.Error(c, "AddError", "")
		return
	}

	DB.Commit()

	echo.Success(c, ApplyBuy, "", "")
}
