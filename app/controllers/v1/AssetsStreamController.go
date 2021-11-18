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
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/validator"
	"strconv"
)

type AssetsStreamController struct {
	BaseController
}

// AssetsStreamHandler 个人资产
func (h *AssetsStreamController) AssetsStreamHandler(c *gin.Context) {
	// 初始化数据模型结构体
	var (
		params requests.AssetsStream
		result []response.UsersWallet
	)
	_ = c.Bind(&params)
	userId, _ := c.Get("user_id")
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	where := cmap.New().Items()
	where[models.Prefix("$_users_wallet.user_id")] = userId
	if len(params.Type) > 0 {
		where[models.Prefix("$_users_wallet.type")] = params.Type
	}
	DB := mysql.DB.Debug()
	DB.Model(models.UsersWallet{}).
		Where(where).
		Joins(models.Prefix("left join $_trading_pair on $_users_wallet.trading_pair_id = $_trading_pair.id")).
		Order(models.Prefix("$_trading_pair.id asc")).
		Find(&result)
	echo.Success(c, result, "", "")
}

// AssetsTypeHandler 资产类型，获取单个币种余额
func (h *AssetsStreamController) AssetsTypeHandler(c *gin.Context) {
	// 初始化数据模型结构体
	var (
		params requests.TradingPair
		result response.UsersWallet
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
	where := cmap.New().Items()
	where["user_id"] = userId   // 用户id
	where["type"] = params.Type // 钱包类型：1现货 2合约
	if len(params.TradingPairId) > 0 {
		where["trading_pair_id"] = params.TradingPairId // 交易对id
	}
	if len(params.TradingPairName) > 0 {
		where["trading_pair_name"] = params.TradingPairName // 交易对名称
	}
	where["type"] = params.Type
	DB := mysql.DB.Debug()
	DB.Model(models.UsersWallet{}).Where(where).Find(&result)
	if result.Id <= 0 {
		echo.Error(c, "CurrencyIsExist", "")
		return
	}
	echo.Success(c, result, "", "")
}

// TransferHandler 划转
func (h *AssetsStreamController) TransferHandler(c *gin.Context) {
	var (
		params                     requests.Transfer    // 请求参数
		UsersWallet1               response.UsersWallet // 查询钱包1数据
		UsersWallet2               response.UsersWallet // 查询钱包2数据
		data1                      models.WalletStream
		data2                      models.WalletStream
		AmountAfter1, AmountAfter2 float64
	)
	_ = c.Bind(&params)
	userId, _ := c.Get("user_id")
	userInfo, _ := c.Get("user")
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	DB := mysql.DB.Debug().Begin()
	// 查询现货账户
	DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("type", "1").
		Where("trading_pair_id", params.TradingPairId).
		Find(&UsersWallet1)
	// 查询合约账户
	DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("type", "2").
		Where("trading_pair_id", params.TradingPairId).
		Find(&UsersWallet2)
	if UsersWallet1.Id <= 0 || UsersWallet2.Id <= 0 {
		echo.Error(c, "UsersWalletIsNotExist", "")
		return
	}

	// 创建钱包流水
	data1.UserId = userId.(int)                                       // 用户id
	data1.Email = userInfo.(map[string]interface{})["email"].(string) // 用户邮箱
	data1.Amount = params.Num                                         // 流转金额
	Num, _ := strconv.ParseFloat(params.Num, 64)
	// 1 从现货账户划转到合约账户  2 从合约账户划转到现货账户
	switch params.Type {
	case "1":
		AmountAfter1 = UsersWallet1.Available - Num
		if AmountAfter1 < 0 {
			echo.Error(c, "InsufficientBalance", "")
			return
		}
		AmountAfter2 = UsersWallet2.Available + Num
		// 支出流水收集
		data1.Way = "2"                                      // 流转方式 1 收入 2 支出
		data1.Type = "3"                                     // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
		data1.TypeDetail = "3"                               // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
		data1.TradingPairId = UsersWallet1.TradingPairId     // 交易对id
		data1.TradingPairName = UsersWallet1.TradingPairName // 交易对名称
		data1.AmountBefore = UsersWallet1.Available          // 流转前的余额
		data1.AmountAfter = AmountAfter1                     // 流转后的余额

		// 收入流水收集
		data2.Way = "1"                                      // 流转方式 1 收入 2 支出
		data2.Type = "3"                                     // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
		data2.TypeDetail = "3"                               // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
		data2.TradingPairId = UsersWallet2.TradingPairId     // 交易对id
		data2.TradingPairName = UsersWallet2.TradingPairName // 交易对名称
		data2.AmountBefore = UsersWallet2.Available          // 流转前的余额
		data2.AmountAfter = AmountAfter2                     // 流转后的余额

		break
	case "2":
		AmountAfter2 = UsersWallet2.Available - Num
		if AmountAfter2 < 0 {
			echo.Error(c, "InsufficientBalance", "")
			return
		}
		AmountAfter1 = UsersWallet1.Available + Num
		// 支出流水收集
		data1.Way = "2"                                      // 流转方式 1 收入 2 支出
		data1.Type = "3"                                     // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
		data1.TypeDetail = "4"                               // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
		data1.TradingPairId = UsersWallet2.TradingPairId     // 交易对id
		data1.TradingPairName = UsersWallet2.TradingPairName // 交易对名称
		data1.AmountBefore = UsersWallet2.Available          // 流转前的余额
		data1.AmountAfter = AmountAfter2                     // 流转后的余额

		// 收入流水收集
		data2.Way = "1"                                      // 流转方式 1 收入 2 支出
		data2.Type = "3"                                     // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
		data2.TypeDetail = "4"                               // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
		data2.TradingPairId = UsersWallet1.TradingPairId     // 交易对id
		data2.TradingPairName = UsersWallet1.TradingPairName // 交易对名称
		data2.AmountBefore = UsersWallet1.Available          // 流转前的余额
		data2.AmountAfter = AmountAfter1                     // 流转后的余额
		break
	default:
		echo.Error(c, "UsersWalletIsNotExist", "")
		return
	}
	// 修改现货钱包余额
	uErr1 := DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("type", "1").
		Where("trading_pair_id", params.TradingPairId).
		Update("available", AmountAfter1).Error
	if uErr1 != nil {
		DB.Rollback()
		echo.Error(c, "OperationFailed", "")
		return
	}
	// 修改合约钱包余额
	uErr2 := DB.Model(models.UsersWallet{}).
		Where("user_id", userId).
		Where("type", "2").
		Where("trading_pair_id", params.TradingPairId).
		Update("available", AmountAfter2).Error
	if uErr2 != nil {
		DB.Rollback()
		echo.Error(c, "OperationFailed", "")
		return
	}

	cErr1 := DB.Model(models.WalletStream{}).Create(&data1).Error
	if cErr1 != nil {
		logger.Error(errors.New(fmt.Sprintf("创建钱包流水失败，%v", cErr1.Error())))
		DB.Rollback()
		echo.Error(c, "LiquidationUnsuccessful", "")
		return
	}

	cErr2 := DB.Model(models.WalletStream{}).Create(&data2).Error
	if cErr2 != nil {
		logger.Error(errors.New(fmt.Sprintf("创建钱包流水失败，%v", cErr2.Error())))
		DB.Rollback()
		echo.Error(c, "LiquidationUnsuccessful", "")
		return
	}
	DB.Commit()
	echo.Success(c, "", "ok", "")
}
