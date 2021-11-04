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
	// 期权合约交易手续费小于零
	if helpers.StringToInt(Currency.FeeOptionContract) < 0 {
		logger.Error(errors.New(fmt.Sprintf("期权合约交易手续费小于零: %v", Currency.FeeOptionContract)))
		echo.Error(c, "FeeOptionContractIsError", "")
		return
	}
	addData.CurrencyName = Currency.Name + "/" + Currency.TradingPairName // 币种名称 例如：BTC/USDT（币种/交易对）
	addData.TradingPairId = helpers.IntToString(Currency.TradingPairId)   // 交易对ID
	addData.TradingPairName = Currency.TradingPairName                    // 交易对名称
	addData.HandleFee = Currency.FeeOptionContract                        // 期权合约交易手续费
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
	// Type 流转类型 1 币币交易 2 永续合约 3 期权合约 4 申购交易 5 划转 6 充值 7 提现 8 冻结
	// TypeDetail 流转详细类型  1 USDT充值  2 银行卡充值  3 币币交易手续费  4 永续合约手续费  5 期权合约手续费  6 币币账户划转到合约账户  7 合约账户划转到币币账户  8 申购冻结  9 币币交易  10 永续合约  11 期权合约
	WalletStream, err4 := new(models.WalletStream).SetAddData("2", "3", "11", addData, Currency, UsersWallet)
	if err4 != nil {
		DB.Rollback()
		echo.Error(c, "AddError", err4.Error())
		return
	}
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
	DB.Commit()
	echo.Success(c, addData, "ok", "")
}
