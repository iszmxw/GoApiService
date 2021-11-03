package routes

import (
	"github.com/gin-gonic/gin"
	controllers "goapi/app/controllers/v1"
	middlewares "goapi/app/middlewares/v1"
)

var v1Group = new(controllers.Group)

// RegisterWebRoutes 注册路由
func RegisterWebRoutes(router *gin.RouterGroup) {
	// 路由分组 客户端 模块
	ApiRoute := router.Group("/api")
	{
		ApiRoute.Use(middlewares.Client())
		// 登录系统
		ApiRoute.Any("/user/login", v1Group.LoginController.LoginHandler)
		// 发送注册邮件
		ApiRoute.Any("/user/send_email_register", v1Group.LoginController.SendEmailRegisterHandler)
		// 找回密码邮件
		ApiRoute.Any("/user/send_email_retrieve", v1Group.LoginController.SendEmailRetrieveHandler)
		// 保存注册信息
		ApiRoute.Any("/user/verify_register", v1Group.LoginController.VerifyRegisterHandler)
		// 重置密码前验证邮件Code
		ApiRoute.Any("/user/reset_verify", v1Group.LoginController.ResetVerifyHandler)
		// 重置密码
		ApiRoute.Any("/user/reset_password", v1Group.LoginController.ResetPasswordHandler)
		// 退出系统
		ApiRoute.Any("/user/logout", v1Group.LoginController.LogoutHandler)
		// 获取用户信息
		ApiRoute.Any("/user/info", v1Group.UserController.UserInfoHandler)
		// 设置语言
		ApiRoute.Any("/user/lang_setup", v1Group.UserController.LangSetupHandler)
		// 设置支付密码
		ApiRoute.Any("/user/payPassword_setup", v1Group.UserController.PayPasswordSetupHandler)
		// 设置支付密码
		ApiRoute.Any("/user/edit_password", v1Group.UserController.EditPasswordHandler)
		// 我的邀请码
		ApiRoute.Any("/user/share_code", v1Group.UserController.ShareCodeHandler)

		// 首页
		index := ApiRoute.Group("/index")
		{
			// 系统信息
			index.Any("/system_info", v1Group.IndexController.SystemInfoHandler)
			// 币种列表
			index.Any("/sys_currency", v1Group.IndexController.SysCurrencyHandler)
			// 交易对列表
			index.Any("/trading_pair", v1Group.IndexController.TradingPairHandler)
		}

		// 资产流水
		assetsStream := ApiRoute.Group("/assets_stream")
		{
			// 个人资产
			assetsStream.Any("/", v1Group.AssetsStreamController.AssetsStreamHandler)
			// 单个币种余额
			assetsStream.Any("/assets_type", v1Group.AssetsStreamController.AssetsTypeHandler)
			// todo::划转暂时不做
			assetsStream.Any("/transfer", v1Group.AssetsStreamController.TransferHandler)
			// 资产流水-订单时间(订单类型/全部交易对)
			assetsStream.Any("/list", v1Group.OrderController.ListHandler)
		}

		// 交易相关
		Trade := ApiRoute.Group("/trade")
		{
			// 获取类型
			Trade.Any("/sys_type", v1Group.TradeController.SsyTypeHandler)
			// 充值
			Trade.Any("/recharge", v1Group.TradeController.ReChargeHandler)
			// 充值记录
			Trade.Any("/recharge_log", v1Group.TradeController.ReChargeLogHandler)
			// 获取提币页余额和提现手续费
			Trade.Any("/get_withdraw_config", v1Group.TradeController.GetWithdrawConfigHandler)
			// 提币
			Trade.Any("/withdraw", v1Group.TradeController.WithdrawHandler)
			// 提币记录
			Trade.Any("/withdraw_log", v1Group.TradeController.WithdrawLogHandler)

			// 申购
			applyBuy := Trade.Group("/apply_buy")
			{
				// 获取申购币种
				applyBuy.Any("/get_currency", v1Group.TradeController.GetCurrencyHandler)
				// 申购
				applyBuy.Any("/submit", v1Group.TradeController.SubmitApplyBuyHandler)
			}
		}

		// 期权合约
		optionContract := ApiRoute.Group("/option_contract")
		{
			// 合约信息列表
			optionContract.Any("/contract_list", v1Group.OptionContractController.ContractListHandler)
			// 订单记录/持单记录
			optionContract.Any("/log", v1Group.OptionContractController.LogHandler)
			// 买张、买跌、自输入
			optionContract.Any("/trade", v1Group.OptionContractController.TradeHandler)
		}

		// 永续合约
		perpetualContract := ApiRoute.Group("/perpetual_contract")
		{
			// 合约信息列表
			perpetualContract.Any("/contract_list", v1Group.PerpetualContractController.ContractListHandler)
			// 历史委托
			perpetualContract.Any("/history", v1Group.PerpetualContractController.HistoryHandler)
			// 限价/市价
			perpetualContract.Any("/trade", v1Group.PerpetualContractController.TradeHandler)
			// 撤单
			perpetualContract.Any("/cancel_order", v1Group.PerpetualContractController.CancelOrderHandler)
			// 平仓
			perpetualContract.Any("/liquidation", v1Group.PerpetualContractController.LiquidationHandler)
		}

		// 币币交易
		currencyCurrency := ApiRoute.Group("/currency_currency")
		{
			// 历史委托
			currencyCurrency.Any("/history", v1Group.CurrencyCurrencyController.HistoryHandler)
			// 买入、买卖
			currencyCurrency.Any("/transaction", v1Group.CurrencyCurrencyController.TransactionHandler)
			// 撤单
			currencyCurrency.Any("/cancel_order", v1Group.CurrencyCurrencyController.CancelOrderHandler)
		}

	}
}
