package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/agent_dividend"
	"goapi/pkg/echo"
	"goapi/pkg/helpers"
	"goapi/pkg/huobi"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/redis"
	"goapi/pkg/validator"
	"strconv"
	"time"
)

type OptionContractController struct {
	BaseController
}

// ContractListHandler 期权合约-合约信息列表
func (h *OptionContractController) ContractListHandler(c *gin.Context) {
	var (
		list []response.OptionContract
	)
	// 查询启用状态的合约
	mysql.DB.Debug().Model(models.OptionContract{}).Where("status", 1).Find(&list)
	if len(list) <= 0 {
		echo.Error(c, "ContractIsNotExist", "")
		return
	}
	echo.Success(c, list, "ok", "")
}

// LogHandler 期权合约-记录
func (h *OptionContractController) LogHandler(c *gin.Context) {
	var (
		request requests.ListOptionContractTransaction // 接收参数
		lists   models.PageList                        // 返回数据
		DB      models.OptionContractTransaction       // 数据模型
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
	if len(request.CurrencyId) > 0 {
		where["currency_id"] = request.CurrencyId
	}
	if len(request.Status) > 0 {
		where["status"] = request.Status
	}
	// 绑定接收的 json 数据到结构体中
	DB.GetPaginate(where, request.OrderBy, int64(request.Page), int64(request.Limit), &lists)
	logger.Info(fmt.Sprintf("%T", lists.Data))
	arr := lists.Data.([]response.OptionContractTransaction)
	for i := 0; i < len(arr); i++ {
		arr[i].Time = arr[i].Seconds
		if arr[i].Status == 0 && len(request.Status) > 0 {
			// 计算倒计时
			arr[i].Seconds = arr[i].CreatedAt.Unix() + arr[i].Seconds - time.Now().Unix()
			// 倒计时小于零，有查询条件的时候
			if arr[i].Seconds <= 0 && len(request.Status) > 0 {
				arr = append(arr[:i], arr[i+1:]...)
				if lists.Total > 0 {
					lists.Total--
				}
				i--
			}
		} else {
			arr = append(arr[:i], arr[i:]...)
		}
	}
	lists.Data = arr
	echo.Success(c, lists, "ok", "")
}

// TradeHandler 期权合约-买张、买跌、自输入
func (h *OptionContractController) TradeHandler(c *gin.Context) {
	var (
		request        requests.OptionContractTransaction
		OptionContract response.OptionContract
		addData        models.OptionContractTransaction
	)
	_ = c.Bind(&request)
	// 数据验证
	vErr := validator.Validate.Struct(request)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, request, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	userId, _ := c.Get("user_id")
	addData.Status = "0" // 状态：0 交易中 1 已完成 2 已撤单
	addData.OrderNumber = helpers.OrderId("OC")
	addData.UserId = userId.(int)
	addData.Email = userInfo.(map[string]interface{})["email"].(string)
	CurrencyId, err := strconv.Atoi(request.CurrencyId)
	if err != nil {
		echo.Error(c, "ValidatorError", err.Error())
		return
	}
	addData.CurrencyId = CurrencyId                     // 币种
	addData.OptionContractId = request.OptionContractId // 期权合约id
	addData.Seconds = request.Seconds                   // 秒数
	addData.ProfitRatio = request.ProfitRatio           // 收益率
	addData.OrderType = request.OrderType               // 订单类型：1-买涨 2-买跌
	addData.Price = request.Price                       // 交易金额
	addData.BuyPrice = request.BuyPrice                 // 购买价格
	// 开启数据库
	DB := mysql.DB.Debug().Begin()
	// 查询交易的合约信息
	DB.Model(models.OptionContract{}).Where("id", request.OptionContractId).Find(&OptionContract)
	if OptionContract.Id <= 0 {
		echo.Error(c, "ContractIsNotExist", "")
		return
	}
	// 期权合约秒数不正确
	if request.Seconds != OptionContract.Seconds {
		echo.Error(c, "OptionContractSecondsNotCorrect", "")
		return
	}
	// 期权合约收益率不正确
	if request.ProfitRatio != OptionContract.ProfitRatio {
		echo.Error(c, "OptionContractProfitRatioNotCorrect", "")
		return
	}
	// 状态:0.禁用,1.启用
	if OptionContract.Status != 1 {
		echo.Error(c, "OptionContractStatus", "")
		return
	}
	// 查询交易的合约信息
	// 查询交易币种信息
	var Currency response.Currency
	DB.Model(models.Currency{}).
		Where(models.Prefix("$_currency.id"), request.CurrencyId).
		Select(models.Prefix("$_currency.*,$_trading_pair.name as trading_pair_name")).
		Joins(models.Prefix("left join $_trading_pair on $_trading_pair.id=$_currency.trading_pair_id")).
		Find(&Currency)
	if Currency.Id <= 0 {
		echo.Error(c, "CurrencyIsExist", "")
		return
	}
	if Currency.DecimalScale > 0 {
		logger.Error(errors.New("自有币种不能进行交易"))
		echo.Error(c, "CurrencyTransactionIsExist", "")
		return
	}
	// 期权合约交易手续费小于零
	if Currency.FeeOptionContract < 0 {
		logger.Error(errors.New(fmt.Sprintf("期权合约交易手续费小于零: %v", Currency.FeeOptionContract)))
		echo.Error(c, "FeeOptionContractIsError", "")
		return
	}
	addData.CurrencyName = Currency.Name + "/" + Currency.TradingPairName // 币种名称 例如：BTC/USDT（币种/交易对）
	addData.TradingPairId = helpers.IntToString(Currency.TradingPairId)   // 交易对ID
	addData.TradingPairName = Currency.TradingPairName                    // 交易对名称
	addData.HandleFee = fmt.Sprintf("%v", Currency.FeeOptionContract)     // 期权合约交易手续费
	// 查询用户钱包信息
	where := cmap.New().Items()
	where["user_id"] = userId
	where["type"] = "2" // 钱包类型：1现货 2合约
	where["trading_pair_id"] = Currency.TradingPairId
	var UsersWallet response.UsersWallet
	DB.Model(models.UsersWallet{}).Where(where).Find(&UsersWallet)
	Price, err2 := strconv.ParseFloat(request.Price, 64)
	if err2 != nil {
		logger.Error(errors.New(fmt.Sprintf("金额转换为float64失败: %v", request.Price)))
		DB.Rollback()
		echo.Error(c, "AddError", err2.Error())
		return
	}
	// 用户可用余额不足
	if UsersWallet.Id <= 0 || UsersWallet.Available <= 0 || UsersWallet.Available < Price {
		logger.Error(errors.New(fmt.Sprintf("UsersWallet.Available: %v <= 0 || UsersWallet.Id: %v <= 0", UsersWallet.Available, UsersWallet.Id)))
		DB.Rollback()
		echo.Error(c, "InsufficientBalance", "")
		return
	}
	if request.OrderType == "1" {
		addData.Remark = "买涨"
	}
	if request.OrderType == "2" {
		addData.Remark = "买跌"
	}
	cErr := DB.Model(&models.OptionContractTransaction{}).Create(&addData).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}

	// 记录资产流水 todo::暂时没有用到后面可能砍掉
	AssetsStream := new(models.AssetsStream).SetAddData(3, addData, Currency)
	cErr = DB.Model(AssetsStream).Create(&AssetsStream).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}

	// 修改钱包余额 （交易扣减）
	UpdateUsersWallet := cmap.New().Items()
	UpdateUsersWallet["available"] = UsersWallet.Available - Price
	editError := DB.Model(models.UsersWallet{}).Where(where).Updates(UpdateUsersWallet).Error
	if editError != nil {
		logger.Error(errors.New(fmt.Sprintf("修改钱包余额失败，%v", editError.Error())))
		DB.Rollback()
		return
	}
	// 修改钱包余额

	// 记录钱包流水
	// Way 流转方式 1 收入 2 支出
	// Type 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	// TypeDetail 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	WalletStream, err4 := new(models.WalletStream).SetAddData("2", "8", "10", addData, Currency, UsersWallet)
	if err4 != nil {
		DB.Rollback()
		echo.Error(c, "AddError", err4.Error())
		return
	}
	WalletStream.HandlingFee = addData.HandleFee // 记录手续费
	cErr = DB.Model(WalletStream).Create(&WalletStream).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}
	ID := helpers.IntToString(addData.Id)
	_, rErr := redis.Add("option_contract:order:"+ID, "期权合约秒交易订单处理", request.Seconds)
	if rErr != nil {
		logger.Error(errors.New("期权合约Redis缓存失败：" + rErr.Error()))
		DB.Rollback()
		echo.Error(c, "AddError", rErr.Error())
		return
	}

	// 检测分发代理分润
	var UserInfo response.User
	DB.Model(models.User{}).Where("id", userId).Find(&UserInfo) // 查询用户信息
	if UserInfo.ParentId > 0 {                                  // parentId 上级代理id
		var arr agent_dividend.Params
		arr.UserId = UserInfo.ParentId                // 用户id
		arr.Email = UserInfo.Email                    // 用户邮箱
		arr.WalletType = 2                            // 钱包类型：1现货 2合约
		arr.TradingPairId = Currency.TradingPairId    // 交易对id
		arr.TradingPairName = addData.TradingPairName // 交易对名称
		arr.TransactionAmount = Price                 // 交易金额
		arr.ParentDividend = 0                        // 上级获得的分润,初始化该值，最底层代理的上级代理分润默认为零，用来后面计算
		arr.WalletStreamType = "8"                    // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
		arr.WalletStreamTypeDetail = "13"             // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
		arr.Current = 10                              // 层级默认只处理10层关系
		// 开启一个 goroutine 去处理
		go agent_dividend.ParentAgentDividend(arr)
	}

	DB.Commit()
	echo.Success(c, addData, "ok", "")
}

// LiquidationHandler 期权合约-平仓
func (h *OptionContractController) LiquidationHandler(c *gin.Context) {
	var params requests.OptionContractTransactionLiquidation
	_ = c.Bind(&params) // 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	var result response.OptionTransactionKline
	DB := mysql.DB.Begin().Debug()
	DB.Model(models.OptionContractTransaction{}).
		Where(models.Prefix("$_option_contract_transaction.id"), params.Id).
		Where(models.Prefix("$_option_contract_transaction.email"), userInfo.(map[string]interface{})["email"]).
		Where(models.Prefix("$_option_contract_transaction.status"), "0").
		Joins(models.Prefix("left join $_currency on $_option_contract_transaction.currency_id=$_currency.id")).
		Select(models.Prefix("$_option_contract_transaction.*,$_currency.k_line_code")).
		Find(&result)
	if result.Id <= 0 {
		DB.Rollback()
		echo.Error(c, "ValidatorError", "id error")
		return
	}
	if result.Status > 0 {
		DB.Rollback()
		logger.Error(errors.New("该订单已确认"))
		echo.Error(c, "ValidatorError", "")
		return
	}

	clinchPrice, err1 := huobi.Kline(result.KLineCode, "close")
	Updates := cmap.New().Items()
	Updates["clinch_price"] = clinchPrice // 成交价格
	if err1 != nil {
		DB.Rollback()
		logger.Error(errors.New(fmt.Sprintf("获取本阶段收盘价失败: %v", err1.Error())))
		echo.Error(c, "OperationFailed", "")
		return
	}
	logger.Info(fmt.Sprintf("获取本阶段收盘价: %v", clinchPrice))
	BuyPrice, err2 := strconv.ParseFloat(result.BuyPrice, 64)
	if err2 != nil {
		DB.Rollback()
		logger.Error(errors.New(fmt.Sprintf("字符串转float64失败: %v", err2.Error())))
		echo.Error(c, "OperationFailed", "")
		return
	}
	// 更新交割结果
	Updates["status"] = "1"
	// 计算交割结果
	// 涨：收盘价大于开盘价
	logger.Info("clinchPrice > BuyPrice")
	if clinchPrice > BuyPrice {
		logger.Info(fmt.Sprintf("%v > %v ：涨", clinchPrice, BuyPrice))
		Updates["order_result"] = "1"
	} else {
		logger.Info(fmt.Sprintf("%v > %v ：跌", clinchPrice, BuyPrice))
		// 跌：收盘价小于开盘价
		Updates["order_result"] = "2"
	}

	// 根据风控来交易，不看k线图结果
	var User response.User
	DB.Model(models.User{}).Where("id", result.UserId).Find(&User)
	// User.RiskProfit 风控 0-无 1-盈 2-亏  根据风控来交易，不看k线图结果   ||    购买结果和实际结果一样 盈利了
	// 盈利：本金+（本金*盈利率）-手续费 handle_fee
	if User.RiskProfit == 2 { // 根据风控概率直接让用户输，篡改交割结果改成和用户买的不一样
		if result.OrderType == Updates["order_result"] {
			if Updates["order_result"] == "1" {
				Updates["order_result"] = "2"
			} else {
				Updates["order_result"] = "1"
			}
			// 最终盈利
			Updates["result_profit"] = 0
		}
	} else {
		if User.RiskProfit == 1 || result.OrderType == Updates["order_result"] {
			if User.RiskProfit == 1 {
				Updates["order_result"] = "1"
				logger.Info(fmt.Sprintf("根据风控概率赢了 User.RiskProfit : %v", User.RiskProfit))
				logger.Info(fmt.Sprintf("风控盈利更新交割结果 : %v", Updates))
			}
			Profit := result.Price + (result.Price * result.ProfitRatio / 100) - (result.Price * result.HandleFee / 100)
			logger.Info(fmt.Sprintf("%v + (%v * %v) - (%v * %v)", result.Price, result.Price, result.ProfitRatio/100, result.Price, result.HandleFee/100))
			logger.Info(fmt.Sprintf("最终盈利 : %v", Profit))
			Updates["result_profit"] = Profit

			// 查询用户当前钱包余额
			where := cmap.New().Items()
			where["user_id"] = result.UserId
			where["trading_pair_id"] = result.TradingPairId
			var UsersWallet response.UsersWallet
			DB.Model(models.UsersWallet{}).Where(where).Find(&UsersWallet)
			// 查询用户当前钱包余额

			// 修改钱包余额
			UpdateUsersWallet := cmap.New().Items()
			UpdateUsersWallet["available"] = UsersWallet.Available + Profit
			editError := DB.Model(models.UsersWallet{}).Where(where).Updates(UpdateUsersWallet).Error
			if editError != nil {
				logger.Error(errors.New(fmt.Sprintf("修改钱包余额失败 : %v", editError.Error())))
				DB.Rollback()
				echo.Error(c, "OperationFailed", "")
				return
			}
			// 修改钱包余额

			// 创建钱包流水
			var data models.WalletStream
			data.Way = "1"                                                 // 流转方式 1 收入 2 支出
			data.Type = "7"                                                // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
			data.TypeDetail = "11"                                         // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
			data.UserId = result.UserId                                    // 用户id
			data.Email = result.Email                                      // 用户邮箱
			data.TradingPairId = helpers.StringToInt(result.TradingPairId) // 交易对id
			data.TradingPairName = result.TradingPairName                  // 交易对名称
			data.Amount = fmt.Sprintf("%v", Profit)                        // 流转金额
			data.AmountBefore = UsersWallet.Available                      // 流转前的余额
			data.AmountAfter = Profit + UsersWallet.Available              // 流转后的余额
			cErr := DB.Model(models.WalletStream{}).Create(&data).Error
			if cErr != nil {
				logger.Error(errors.New(fmt.Sprintf("创建钱包流水失败 : %v", cErr.Error())))
				DB.Rollback()
				echo.Error(c, "OperationFailed", "")
				return
			}
			// 创建钱包流水
		}
		logger.Info("clinchPrice > BuyPrice")
	}
	err3 := DB.Model(models.OptionContractTransaction{}).
		Where("id", params.Id).
		Updates(Updates).Error
	if err3 != nil {
		logger.Error(errors.New(fmt.Sprintf("更新交割结果失败 : %v", err3.Error())))
		DB.Rollback()
		echo.Error(c, "OperationFailed", "")
		return
	}
	DB.Commit()
	echo.Success(c, "", "", "")

}
