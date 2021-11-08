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
	where[models.Prefix("$_currency_transaction.email")] = email
	if len(request.CurrencyId) > 0 {
		where[models.Prefix("$_currency_transaction.currency_id")] = request.CurrencyId
	}
	if len(request.Status) > 0 {
		where[models.Prefix("$_currency_transaction.status")] = request.Status
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
	addData.EntrustNum = params.EntrustNum           // 委托量
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

	// 1、查询用户钱包对的钱包信息
	where := cmap.New().Items()
	where["user_id"] = userId
	where["status"] = "0" // 0正常 1锁定
	where["type"] = "1"   // 钱包类型：1现货 2合约
	where["trading_pair_id"] = Currency.TradingPairId
	var UsersWallet response.UsersWallet
	DB.Model(models.UsersWallet{}).Where(where).Find(&UsersWallet)
	// 用户钱包对可用余额不足
	if UsersWallet.Available <= 0 || UsersWallet.Id <= 0 {
		log := fmt.Sprintf("UsersWallet.Available（%v） <= 0 || UsersWallet.Id（%v） <= 0", UsersWallet.Available, UsersWallet.Id)
		logger.Info(UsersWallet)
		logger.Error(errors.New(log))
		DB.Rollback()
		echo.Error(c, "InsufficientBalance", "")
		return
	}
	// 2、查询用户交易对的钱包信息
	where2 := cmap.New().Items()
	where2["user_id"] = userId
	where2["status"] = "0" // 0正常 1锁定
	where2["type"] = "1"   // 钱包类型：1现货 2合约
	where2["trading_pair_name"] = Currency.Name
	var UsersWallet2 response.UsersWallet
	DB.Model(models.UsersWallet{}).Where(where2).Find(&UsersWallet2)
	// 交易对钱包信息不存在，或者可用余额不足
	if UsersWallet2.Available <= 0 || UsersWallet2.Id <= 0 {
		log := fmt.Sprintf("UsersWallet2.Available（%v） <= 0 || UsersWallet2.Id（%v） <= 0", UsersWallet2.Available, UsersWallet2.Id)
		logger.Info(UsersWallet2)
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
	var WalletStreamUsersWallet response.UsersWallet
	UpdateUsersWallet := cmap.New().Items()
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
		// 可用余额减去消费货币
		UpdateUsersWallet["available"] = UsersWallet.Available - buyNum
		WalletStreamUsersWallet = UsersWallet
		// 订单方向：1-买入 // 修改钱包余额 （交易扣减）
		editError := DB.Model(models.UsersWallet{}).Where(where).Updates(UpdateUsersWallet).Error
		if editError != nil {
			logger.Info("修改钱包余额失败")
			logger.Info(editError.Error())
			DB.Rollback()
			return
		}
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
		// 可用余额减去卖出货币
		UpdateUsersWallet["available"] = UsersWallet2.Available - Percentage
		WalletStreamUsersWallet = UsersWallet2
		// 订单方向：2-卖出 // 修改钱包余额 （交易扣减）
		editError := DB.Model(models.UsersWallet{}).Where(where2).Updates(UpdateUsersWallet).Error
		if editError != nil {
			logger.Info("修改钱包余额失败")
			logger.Info(editError.Error())
			DB.Rollback()
			return
		}
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

	// 记录资产流水 todo 暂未用到后面可能砍掉
	AssetsStream := new(models.AssetsStream).SetAddData(1, addData, Currency)
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
	WalletStream, err4 := new(models.WalletStream).SetAddData("2", "6", "8", addData, Currency, WalletStreamUsersWallet)
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
		Select(models.Prefix("$_currency_transaction.*,$_currency.trading_pair_id,$_currency.name")).Find(&CurrencyTransaction)
	if CurrencyTransaction.Status == "0" {
		// 买入订单撤销
		if CurrencyTransaction.TransactionType == "1" {
			// 钱包搜索条件
			whereUsersWallet := map[string]interface{}{"user_id": userId, "type": "1", "trading_pair_id": CurrencyTransaction.TradingPairId}
			DB.Model(models.UsersWallet{}).Where(whereUsersWallet).Find(&UsersWallet)
			// 退还消费的金额到账户
			UsersWallet.Available = UsersWallet.Available + CurrencyTransaction.Price
			uErr := DB.Model(models.UsersWallet{}).Where(whereUsersWallet).Update("available", UsersWallet.Available).Error
			if uErr != nil {
				fmt.Println(uErr.Error())
				DB.Rollback()
				echo.Error(c, "OperationFailed", uErr.Error())
				return
			}
		}

		// 卖出订单撤销
		if CurrencyTransaction.TransactionType == "2" {
			// 钱包搜索条件
			whereUsersWallet := map[string]interface{}{"user_id": userId, "type": "1", "trading_pair_name": CurrencyTransaction.Name}
			DB.Model(models.UsersWallet{}).Where(whereUsersWallet).Find(&UsersWallet)
			// 退还卖出的金额到账户
			ClinchNum, _ := strconv.ParseFloat(CurrencyTransaction.ClinchNum, 64)
			UsersWallet.Available = UsersWallet.Available + ClinchNum
			uErr := DB.Model(models.UsersWallet{}).Where(whereUsersWallet).Update("available", UsersWallet.Available).Error
			if uErr != nil {
				fmt.Println(uErr.Error())
				DB.Rollback()
				echo.Error(c, "OperationFailed", uErr.Error())
				return
			}
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
