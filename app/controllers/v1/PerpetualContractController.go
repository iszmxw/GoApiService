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
	"goapi/pkg/validator"
	"strconv"
	"strings"
)

type PerpetualContractController struct {
	BaseController
}

// HistoryHandler 永续合约-历史委托
func (h *PerpetualContractController) HistoryHandler(c *gin.Context) {
	var (
		request requests.ListPerpetualContractTransaction // 接收参数
		lists   models.PageList                           // 返回数据
		DB      models.PerpetualContractTransaction       // 数据模型
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
	echo.Success(c, lists, "ok", "")
}

// ContractListHandler 永续合约-合约信息列表
func (h *PerpetualContractController) ContractListHandler(c *gin.Context) {
	var (
		params requests.PerpetualContract
		list   []response.PerpetualContract
	)
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	where := cmap.New().Items()
	where["currency_id"] = params.CurrencyId
	// 查询启用状态的合约
	mysql.DB.Debug().Model(models.PerpetualContract{}).Where(where).Find(&list)
	if len(list) <= 0 {
		echo.Error(c, "ContractIsNotExist", "")
		return
	}
	echo.Success(c, list, "ok", "")
}

// TradeHandler 永续合约-限价/市价
func (h *PerpetualContractController) TradeHandler(c *gin.Context) {
	var (
		params            requests.PerpetualContractTransaction // 绑定接收的 json 数据到结构体中
		PerpetualContract response.PerpetualContract            // 查询永续合约信息
		Currency          response.Currency                     // 查询交易币种信息
		UsersWallet       response.UsersWallet                  // 查询钱包信息
		UserStatus        response.User                         // 查询用户状态
		addData           models.PerpetualContractTransaction   // 添加永续合约交易数据
		Bail              string                                // 保证金占用比例Bail
		EnsureAmount      float64                               // 保证金
	)
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	// 根据当前交易币种检查交易的合约信息
	mysql.DB.Debug().Model(models.PerpetualContract{}).Where("currency_id", params.CurrencyId).Find(&PerpetualContract)
	if PerpetualContract.Id <= 0 {
		echo.Error(c, "ContractIsNotExist", "")
		return
	}
	arrayMultiple := strings.Split(PerpetualContract.Multiple, ",") // 倍数
	arrayBail := strings.Split(PerpetualContract.Bail, ",")         // 保证金占用比例
	// 检查合约对应的信息是否设置正确
	if len(arrayMultiple) != len(arrayBail) {
		echo.Error(c, "ContractIsNotCorrect", "")
		return
	}
	// 通过倍数值获取相对应的保证金占比值
	for key, val := range arrayMultiple {
		if params.Multiple == val {
			Bail = arrayBail[key]
		}
	}
	// 该函数会打乱数组原有的顺序
	if !helpers.InArray(params.Multiple, arrayMultiple) {
		echo.Error(c, "MultipleIsError", "")
		return
	}
	logger.Info(fmt.Sprintf("获取保证金占用比例Bail:%v", Bail))
	// 根据当前交易币种查询交易的合约信息

	// 接收的字符串参数转换float64处理
	EntrustNum, EntrustNumErr := strconv.ParseFloat(params.EntrustNum, 64) // 委托量，委托（20%，50%，75%，100%）
	LimitPrice, LimitPriceErr := strconv.ParseFloat(params.LimitPrice, 64) // todo::当前限价(限价的时候手输入，市价的时候，传入k线图的最高价，后期后台自动去火币网获取)
	Bails, BailsErr := strconv.ParseFloat(Bail, 64)                        // 保证金占比值
	CurrencyId, CurrencyIdErr := strconv.Atoi(params.CurrencyId)           // 币种id
	if EntrustNumErr != nil {
		echo.Error(c, "EntrustNumErr", "")
		return
	}
	if LimitPriceErr != nil {
		echo.Error(c, "LimitPriceErr", "")
		return
	}
	if BailsErr != nil {
		echo.Error(c, "Percentage", "")
		return
	}
	if CurrencyIdErr != nil {
		echo.Error(c, "ValidatorError", CurrencyIdErr.Error())
		return
	}

	userInfo, _ := c.Get("user")
	userId, _ := c.Get("user_id")
	addData.Status = "0"                        // 状态：0 交易中 1 已完成 2 已撤单
	addData.OrderNumber = helpers.OrderId("PC") // 订单号
	addData.Email = userInfo.(map[string]interface{})["email"].(string)
	addData.CurrencyId = CurrencyId            // 币种
	addData.UserId = userId.(int)              // 用户id
	addData.EntrustNum = params.EntrustNum     // 委托量
	addData.EntrustPrice = params.EntrustPrice // 委托价格
	addData.LimitPrice = params.LimitPrice     // 当前限价
	// 保证金
	EnsureAmount = EntrustNum * LimitPrice * (Bails / 100)
	logger.Info(fmt.Sprintf("最终保证金：%v", fmt.Sprintf("%.8f", EnsureAmount)))
	addData.EnsureAmount = fmt.Sprintf("%.8f", EnsureAmount) // 保证金
	addData.Multiple = params.Multiple                       // 倍数值
	addData.OrderType = params.OrderType                     // 订单类型：1-限价 2-市价
	addData.TransactionType = params.TransactionType         // 交易类型：1-开多 2-开空

	// 卖出所得货币=市价（限价）*卖出数量
	logger.Info(fmt.Sprintf("市价/限价: %v * 卖出数量: %v \n", LimitPrice, EntrustNum))
	Price := LimitPrice * EntrustNum
	addData.Price = fmt.Sprintf("%.8f", Price)
	if len(addData.Price) <= 0 {
		logger.Error(errors.New("交易金额必须大于零"))
		echo.Error(c, "AddError", "")
		return
	}

	// 开启数据库
	DB := mysql.DB.Debug().Begin()
	if UserStatus.Status == "1" {
		DB.Rollback()
		echo.Error(c, "UserIsLock", "")
		return
	}
	DB.Model(models.Currency{}).
		Where(models.Prefix("$_currency.id"), params.CurrencyId).
		Select(models.Prefix("$_currency.*,$_trading_pair.name as trading_pair_name")).
		Joins(models.Prefix("left join $_trading_pair on $_trading_pair.id=$_currency.trading_pair_id")).
		Find(&Currency)
	if Currency.Id <= 0 {
		echo.Error(c, "CurrencyIsExist", "")
		return
	}
	arrayType := strings.Split(Currency.Type, ",")
	logger.Info(arrayType)
	// 该函数会打乱数组原有的顺序
	if !helpers.InArray("2", arrayType) {
		DB.Rollback()
		echo.Error(c, "CurrencyTypeIsNotAllowed", "")
		return
	}
	if Currency.DecimalScale > 0 {
		logger.Info(Currency.DecimalScale)
		logger.Error(errors.New("自有币种不能进行交易"))
		echo.Error(c, "CurrencyTransactionIsExist", "")
		return
	}
	// 期权合约交易手续费小于零
	if Currency.FeePerpetualContract < 0 {
		logger.Error(errors.New(fmt.Sprintf("永续合约交易手续费小于零: %v", Currency.FeePerpetualContract)))
		echo.Error(c, "FeePerpetualContractIsError", "")
		return
	}
	addData.CurrencyName = Currency.Name + "/" + Currency.TradingPairName // 币种名称 例如：BTC/USDT（币种/交易对）
	addData.KLineCode = Currency.KLineCode                                // K线图代码
	addData.TradingPairId = Currency.TradingPairId                        // 交易对id
	addData.TradingPairName = Currency.TradingPairName                    // 交易对名称
	addData.HandleFee = fmt.Sprintf("%v", Currency.FeePerpetualContract)  // 手续费百分比
	// 查询用户钱包信息
	where := cmap.New().Items()
	where["user_id"] = userId
	where["type"] = "2" // 钱包类型：1现货 2合约
	where["trading_pair_id"] = Currency.TradingPairId
	DB.Model(models.UsersWallet{}).Where(where).Find(&UsersWallet)
	if UsersWallet.Status == 1 { // 0正常 1锁定
		DB.Rollback()
		echo.Error(c, "UsersWalletLock", "")
		return
	}
	if UsersWallet.Available <= 0 || UsersWallet.Id <= 0 || UsersWallet.Available < Price {
		logger.Info(fmt.Sprintf("UsersWallet.Available: %v <= 0 || UsersWallet.Id: %v <= 0 || UsersWallet.Available:%v < Price:%v", UsersWallet.Available, UsersWallet.Id, UsersWallet.Available, Price))
		DB.Rollback()
		echo.Error(c, "InsufficientBalance", "")
		return
	}
	// 扣钱钱包余额
	UsersWallet.Available = UsersWallet.Available - Price
	vErr = DB.Model(models.UsersWallet{}).Where(where).Update("available", UsersWallet.Available).Error
	if vErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", vErr.Error())
		return
	}
	// 扣钱钱包余额

	// 交易类型 1-开多 2-开空
	if params.TransactionType == "1" {
		addData.Remark = "开多" // 备注
	}
	if params.TransactionType == "2" {
		addData.Remark = "开空" // 备注
	}
	cErr := DB.Model(addData).Create(&addData).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}

	// 记录资产流水 todo::暂时用不到后面可能会砍掉
	AssetsStream := new(models.AssetsStream).SetAddData(2, addData, Currency)
	cErr = DB.Model(AssetsStream).Create(&AssetsStream).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}

	// 记录钱包流水
	// Way 流转方式 1 收入 2 支出
	// Type 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	// TypeDetail 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	WalletStream, err4 := new(models.WalletStream).SetAddData("2", "7", "10", addData, Currency, UsersWallet)
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
	DB.Commit()
	echo.Success(c, addData, "ok", "")
}

// CancelOrderHandler 永续合约-撤单
func (h *PerpetualContractController) CancelOrderHandler(c *gin.Context) {
	var params requests.CancelOrder
	_ = c.Bind(&params) // 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	where := cmap.New().Items()
	where["email"] = userInfo.(map[string]interface{})["email"]
	where["id"] = params.Id
	DB := mysql.DB.Debug()
	uErr := DB.Model(models.PerpetualContractTransaction{}).Where(where).Update("status", "2").Error
	if uErr != nil {
		DB.Rollback()
		echo.Error(c, "OperationFailed", uErr.Error())
		return
	}
	DB.Commit()
	echo.Success(c, "", "", "")
}

// LiquidationHandler 永续合约-平仓
func (h *PerpetualContractController) LiquidationHandler(c *gin.Context) {
	var (
		params                       requests.Liquidation                  // 接收永续合约平仓
		PerpetualContractTransaction response.PerpetualContractTransaction // 永续合约交易信息
		UsersWallet                  response.UsersWallet                  // 钱包列表
		data                         models.WalletStream                   // 钱包流水
		UserInfo                     response.User                         // 查询用户信息
		arr                          agent_dividend.Params                 // 代理信息
		income                       float64                               // 最终收益
	)
	_ = c.Bind(&params) // 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	userId, _ := c.Get("user_id")
	where := cmap.New().Items()
	where["email"] = userInfo.(map[string]interface{})["email"]
	where["id"] = params.Id
	where["status"] = "0"
	DB := mysql.DB.Begin().Debug()
	DB.Model(models.PerpetualContractTransaction{}).Where(where).Find(&PerpetualContractTransaction)
	if PerpetualContractTransaction.Id <= 0 {
		echo.Error(c, "ValidatorError", "id error")
		return
	}
	if PerpetualContractTransaction.Status > 0 {
		logger.Error(errors.New("该订单已确认"))
		echo.Error(c, "ValidatorError", "")
		return
	}
	// 查询用户钱包信息
	whereWallet := cmap.New().Items()
	whereWallet["user_id"] = PerpetualContractTransaction.UserId
	whereWallet["trading_pair_id"] = PerpetualContractTransaction.TradingPairId
	whereWallet["type"] = "2" // 钱包类型：1现货 2合约
	DB.Model(models.UsersWallet{}).Where(whereWallet).Find(&UsersWallet)
	// 用户可用余额不足
	if UsersWallet.Available < 0 || UsersWallet.Id <= 0 {
		logger.Error(errors.New(fmt.Sprintf("UsersWallet.Available: %v < 0 || UsersWallet.Id: %v <= 0", UsersWallet.Available, UsersWallet.Id)))
		DB.Rollback()
		echo.Error(c, "InsufficientBalance", "")
		return
	}

	Liquidation, LiquidationErr := huobi.Kline(PerpetualContractTransaction.KLineCode, "close")
	if LiquidationErr != nil {
		DB.Rollback()
		echo.Error(c, "OperationFailed", "")
		return
	}

	/**
		永续合约占用保证金比例
	  	当20X=1：100
	    50X=1：50
	    100x=1：25
	    200X=1：5
	    // （2.0暂无）可开张数： 可用金额/占用保证金-手续费（手数*手续费百分比）
	*/
	// 平仓时候的k线图，前端传递过来
	// 最终收益 = k线图收盘价 - 委托价格 + 保证金 - 手续费
	if Liquidation < PerpetualContractTransaction.EntrustPrice {
		// 亏损情况
		income = PerpetualContractTransaction.EnsureAmount - (PerpetualContractTransaction.EnsureAmount * PerpetualContractTransaction.HandleFee / 100)
		tips := fmt.Sprintf("最终收益：params.Liquidation - PerpetualContractTransaction.EntrustPrice + PerpetualContractTransaction.EnsureAmount - (PerpetualContractTransaction.EnsureAmount * PerpetualContractTransaction.HandleFee / 100) = income")
		// income最终收益：55499.23 - 55499.01 + 20 - (20 * 0.05 / 100) = 20.210000000001163
		tips += fmt.Sprintf("最终收益：%v - %v + %v - (%v * %v / 100) = %v", Liquidation, PerpetualContractTransaction.EntrustPrice, PerpetualContractTransaction.EnsureAmount, PerpetualContractTransaction.EnsureAmount, PerpetualContractTransaction.HandleFee, income)
		fmt.Println(tips)
	} else {
		// 盈利情况
		income = Liquidation - PerpetualContractTransaction.EntrustPrice + PerpetualContractTransaction.EnsureAmount - (PerpetualContractTransaction.EnsureAmount * PerpetualContractTransaction.HandleFee / 100)
	}
	fmt.Println("k线图代码：", PerpetualContractTransaction.KLineCode)
	fmt.Println("最终盈利：", income)
	//clinchPrice, err1 := huobi.Kline(PerpetualContractTransaction.KLineCode)
	Updates := cmap.New().Items()
	Updates["income"] = income            // 最终收益
	Updates["status"] = "1"               // 确认状态
	Updates["clinch_price"] = Liquidation // 成交价格
	err2 := DB.Model(models.PerpetualContractTransaction{}).Where(where).Updates(Updates).Error
	if err2 != nil {
		logger.Error(errors.New("平仓失败"))
		logger.Error(err2)
		echo.Error(c, "LiquidationUnsuccessful", "")
		return
	}

	// 修改钱包余额 （交易盈利）
	UpdateUsersWallet := cmap.New().Items()
	UpdateUsersWallet["available"] = UsersWallet.Available + income // 返回保证金
	editError := DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("type", "2"). // 钱包类型：1现货 2合约
		Where("trading_pair_id", PerpetualContractTransaction.TradingPairId).
		Updates(UpdateUsersWallet).Error
	if editError != nil {
		logger.Error(errors.New(fmt.Sprintf("修改钱包余额失败，%v", editError.Error())))
		DB.Rollback()
		echo.Error(c, "LiquidationUnsuccessful", "")
		return
	}
	// 修改钱包余额

	// 创建钱包流水
	data.Way = "1"                                                      // 流转方式 1 收入 2 支出
	data.Type = "2"                                                     // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	data.TypeDetail = "10"                                              // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	data.UserId = PerpetualContractTransaction.UserId                   // 用户id
	data.Email = PerpetualContractTransaction.Email                     // 用户邮箱
	data.TradingPairId = PerpetualContractTransaction.TradingPairId     // 交易对id
	data.TradingPairName = PerpetualContractTransaction.TradingPairName // 交易对名称
	data.Amount = fmt.Sprintf("%v", income)                             // 流转金额
	data.AmountBefore = UsersWallet.Available                           // 流转前的余额
	data.AmountAfter = income + UsersWallet.Available                   // 流转后的余额
	cErr := DB.Model(models.WalletStream{}).Create(&data).Error
	if cErr != nil {
		logger.Error(errors.New(fmt.Sprintf("创建钱包流水失败，%v", cErr.Error())))
		DB.Rollback()
		echo.Error(c, "LiquidationUnsuccessful", "")
		return
	}
	// 创建钱包流水

	// 检测分发代理分润
	DB.Model(models.User{}).Where("id", userId).Find(&UserInfo) // 查询用户信息
	if UserInfo.ParentId > 0 {                                  // parentId 上级代理id
		arr.UserId = UserInfo.ParentId                                     // 用户id
		arr.Email = UserInfo.Email                                         // 用户邮箱
		arr.WalletType = 2                                                 // 钱包类型：1现货 2合约
		arr.TradingPairId = PerpetualContractTransaction.TradingPairId     // 交易对id
		arr.TradingPairName = PerpetualContractTransaction.TradingPairName // 交易对名称
		arr.TransactionAmount = PerpetualContractTransaction.Price         // 交易金额
		arr.ParentDividend = 0                                             // 上级获得的分润,初始化该值，最底层代理的上级代理分润默认为零，用来后面计算
		arr.WalletStreamType = "7"                                         // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
		arr.WalletStreamTypeDetail = "11"                                  // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
		arr.Current = 10                                                   // 层级默认只处理10层关系
		// 开启一个 goroutine 去处理
		go agent_dividend.ParentAgentDividend(arr)
	}
	DB.Commit()
	echo.Success(c, "", "", "")
}
