package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/echo"
	"goapi/pkg/helpers"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/validator"
	"strconv"
)

type CurrencyCurrencyController struct {
	BaseController
}

// HistoryHandler 币币交易-记录
func (h *CurrencyCurrencyController) HistoryHandler(c *gin.Context) {
	var (
		request requests.History           // 接收参数
		lists   models.PageList            // 返回数据
		DB      models.CurrencyTransaction // 数据模型
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
	echo.Success(c, lists, "", "")
}

// TransactionHandler 币币交易-买入、卖出   addData.TransactionType = "1" // 订单方向：1-买入 2-卖出
func (h *CurrencyCurrencyController) TransactionHandler(c *gin.Context) {
	var params requests.CurrencyTransaction // 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	userId, _ := c.Get("user_id")
	var addData models.CurrencyTransaction
	addData.Status = "0" // 状态：0 交易中 1 已完成 2 已撤单
	addData.OrderNumber = helpers.OrderId("CC")
	addData.Email = userInfo.(map[string]interface{})["email"].(string)
	CurrencyId, err := strconv.Atoi(params.CurrencyId)
	if err != nil {
		logger.Error(errors.New("传输的币种id转义失败"))
		echo.Error(c, "ValidatorError", err.Error())
		return
	}
	addData.CurrencyId = CurrencyId // 币种
	addData.UserId = userId.(int)   // 用户id
	//addData.EntrustNum = params.EntrustNum           // todo::收集无用--委托量
	addData.LimitPrice = params.LimitPrice           // 当前限价
	addData.ClinchNum = params.ClinchNum             // 成交量
	addData.OrderType = params.OrderType             // 挂单类型：1-限价 2-市价
	addData.TransactionType = params.TransactionType // 订单方向：1-买入 2-卖出
	// 开启数据库
	DB := mysql.DB.Debug().Begin()

	// 查询交易币种信息
	var Currency response.Currency
	DB.Model(models.Currency{}).
		Where(map[string]interface{}{models.Prefix("$_currency.id"): params.CurrencyId}).
		Select(models.Prefix(models.Prefix("$_currency.*,$_trading_pair.name as trading_pair_name"))).
		Joins(models.Prefix("left join $_trading_pair on $_trading_pair.id=$_currency.trading_pair_id")).
		Find(&Currency)
	if Currency.Id <= 0 {
		logger.Error(errors.New("币种不存在"))
		echo.Error(c, "CurrencyIsExist", "")
		return
	}
	addData.CurrencyName = Currency.Name + "/" + Currency.TradingPairName // 币种名称 例如：BTC/USDT（币种/交易对）

	// 查询用户钱包信息
	where := cmap.New().Items()
	where["user_id"] = userId
	where["status"] = "0" // 0正常 1锁定
	where["type"] = "1"   // 钱包类型：1现货 2合约
	where["trading_pair_id"] = Currency.TradingPairId
	var UsersWallet response.UsersWallet
	DB.Model(models.UsersWallet{}).Where(where).Find(&UsersWallet)
	// 用户可用余额不足
	if UsersWallet.Available <= 0 || UsersWallet.Id <= 0 {
		log := fmt.Sprintf("%v <= 0 || %v <= 0", UsersWallet.Available, UsersWallet.Id)
		logger.Info(UsersWallet)
		logger.Error(errors.New(log))
		DB.Rollback()
		echo.Error(c, "InsufficientBalance", "")
		return
	}
	// 限价百分比算法：（当前可用余额*百分比）/当前限价=买入货币

	// 挂单类型：1-限价 2-市价 类型计算
	if len(params.LimitPrice) <= 0 {
		echo.Error(c, "LimitPrice", "")
		return
	}

	var buyNum float64
	// 订单方向：1-买入
	if params.TransactionType == "1" {
		// todo::当前限价(限价的时候手输入，市价的时候，传入k线图的最高价，后期后台自动去火币网获取)
		LimitPrice, err2 := strconv.ParseFloat(params.LimitPrice, 64)
		// 成交量
		Percentage, err1 := strconv.ParseFloat(params.ClinchNum, 64)
		if err1 != nil || err2 != nil {
			echo.Error(c, "Percentage", "")
			return
		}
		// 买入消费货币=市价（限价）*卖出数量
		logger.Info(fmt.Sprintf("市价/限价: %v * 卖出数量: %v \n", LimitPrice, Percentage))
		buyNum = LimitPrice * Percentage
		// 消费的金额不能大于钱包余额
		if buyNum > UsersWallet.Available {
			echo.Error(c, "InsufficientBalance", "")
			return
		}
		addData.Price = fmt.Sprintf("%.8f", buyNum)
		// 挂单类型：1-限价 2-市价 类型计算
	}

	// 订单方向：2-卖出
	if params.TransactionType == "2" {
		// todo::当前限价(限价的时候手输入，市价的时候，传入k线图的最高价，后期后台自动去火币网获取)
		LimitPrice, err2 := strconv.ParseFloat(params.LimitPrice, 64)
		// 成交量
		Percentage, err1 := strconv.ParseFloat(params.ClinchNum, 64)
		if err1 != nil || err2 != nil {
			echo.Error(c, "Percentage", "")
			return
		}
		// 卖出所得货币=市价（限价）*卖出数量
		fmt.Printf("市价/限价: %v * 卖出数量: %v \n", LimitPrice, Percentage)
		buyNum = LimitPrice * Percentage
		addData.Price = fmt.Sprintf("%.8f", buyNum)
		// 挂单类型：1-限价 2-市价 类型计算
	}
	if len(addData.Price) <= 0 {
		logger.Error(errors.New("买入价格计算错误"))
		echo.Error(c, "AddError", "")
		return
	}
	// 计算挂单价格
	if params.TransactionType == "1" {
		addData.Remark = "买入" // 备注
	}
	if params.TransactionType == "2" {
		addData.Remark = "卖出" // 备注
	}
	cErr := DB.Model(addData).Create(&addData).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}

	// 记录资产流水
	AssetsStream := new(models.AssetsStream).SetAddData(1, addData, Currency)
	cErr = DB.Model(AssetsStream).Create(&AssetsStream).Error
	if cErr != nil {
		DB.Rollback()
		echo.Error(c, "AddError", cErr.Error())
		return
	}

	// 修改钱包余额 （交易扣减）
	UpdateUsersWallet := cmap.New().Items()
	UpdateUsersWallet["available"] = UsersWallet.Available - buyNum
	editError := DB.Model(models.UsersWallet{}).Where(where).Updates(UpdateUsersWallet).Error
	if editError != nil {
		logger.Info("修改钱包余额失败")
		logger.Info(editError.Error())
		DB.Rollback()
		return
	}
	// 修改钱包余额

	// 记录钱包流水
	// Way 流转方式 1 收入 2 支出
	// Type 流转类型 1 币币交易 2 永续合约 3 期权合约 4 申购交易 5 划转 6 充值 7 提现 8 冻结
	// TypeDetail 流转详细类型  1 USDT充值  2 银行卡充值  3 币币交易手续费  4 永续合约手续费  5 期权合约手续费  6 币币账户划转到合约账户  7 合约账户划转到币币账户  8 申购冻结  9 币币交易  10 永续合约  11 期权合约
	WalletStream, err4 := new(models.WalletStream).SetAddData("2", "1", "9", addData, Currency, UsersWallet)
	if err4 != nil {
		DB.Rollback()
		logger.Error(errors.New("添加数据失败" + err4.Error()))
		echo.Error(c, "AddError", err4.Error())
		return
	}
	cErr = DB.Model(WalletStream).Create(&WalletStream).Error
	if cErr != nil {
		DB.Rollback()
		logger.Error(errors.New("添加数据失败" + cErr.Error()))
		echo.Error(c, "AddError", cErr.Error())
		return
	}
	DB.Commit()
	echo.Success(c, addData, "", "")
}

// CancelOrderHandler 币币交易-撤单
func (h *CurrencyCurrencyController) CancelOrderHandler(c *gin.Context) {
	var params requests.CancelOrder
	var CurrencyTransaction response.CurrencyTransaction
	var UsersWallet response.UsersWallet
	_ = c.Bind(&params) // 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	userId, _ := c.Get("user_id")
	DB := mysql.DB.Debug()
	DB.Model(models.CurrencyTransaction{}).
		Where(map[string]interface{}{
			models.Prefix("$_currency_transaction.email"): userInfo.(map[string]interface{})["email"],
			models.Prefix("$_currency_transaction.id"):    params.Id},
		).
		Joins(models.Prefix("left join $_currency on $_currency.id=$_currency_transaction.currency_id")).
		Select(models.Prefix("$_currency_transaction.*,$_currency.trading_pair_id")).Find(&CurrencyTransaction)
	if CurrencyTransaction.Status == "0" {
		// 钱包搜索条件
		whereUsersWallet := map[string]interface{}{"user_id": userId, "trading_pair_id": CurrencyTransaction.CurrencyId}
		DB.Model(models.UsersWallet{}).Where(whereUsersWallet).Find(&UsersWallet)
		// 退还金额到账户
		UsersWallet.Available = UsersWallet.Available + CurrencyTransaction.Price
		uErr := DB.Model(models.UsersWallet{}).Where(whereUsersWallet).Update("available", UsersWallet.Available).Error
		if uErr != nil {
			fmt.Println(uErr.Error())
			DB.Rollback()
			echo.Error(c, "OperationFailed", uErr.Error())
			return
		}
	}
	where := cmap.New().Items()
	where["email"] = userInfo.(map[string]interface{})["email"]
	where["id"] = params.Id
	where["status"] = "0"
	uErr := DB.Model(models.CurrencyTransaction{}).Where(where).Update("status", "2").Error
	if uErr != nil {
		DB.Rollback()
		echo.Error(c, "OperationFailed", uErr.Error())
		return
	}
	DB.Commit()
	echo.Success(c, "", "", "")
}
