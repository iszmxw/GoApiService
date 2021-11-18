package agent_dividend

import (
	"errors"
	"fmt"
	"goapi/app/models"
	"goapi/app/response"
	"goapi/pkg/helpers"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"gorm.io/gorm"
	"strconv"
)

// 用户信息表

type ParentAgent struct {
	Id            int                `json:"id"`             //主键id
	Language      string             `json:"language"`       //语言：1-繁体 2-英文 3-日文 4-韩语 5 西班牙语
	IsAgent       string             `json:"is_agent"`       //是否代理： 0：不是   1：是
	ParentId      int                `json:"parent_id"`      //上级用户id
	Email         string             `json:"email"`          //邮箱
	Nickname      string             `json:"nickname"`       //昵称
	Password      string             `json:"password"`       //登录密码
	PayPassword   string             `json:"pay_password"`   //支付密码
	UserLevel     int                `json:"user_level"`     //用户层级
	UserPath      string             `json:"user_path"`      //用户关系
	AgentDividend string             `json:"agent_dividend"` //代理红利
	ShareCode     string             `json:"share_code"`     //用户邀请码，每个用户唯一
	RiskProfit    int                `json:"risk_profit"`    //风控 0-无 1-盈 2-亏
	LastLoginIp   string             `json:"last_login_ip"`  //登录IP
	Status        string             `json:"status"`         //状态： 0正常 1，已锁定
	LockTime      helpers.TimeNormal `json:"lock_time"`      //锁定时间
	LoginTime     helpers.TimeNormal `json:"login_time"`     //登录时间
	CreatedAt     helpers.TimeNormal `json:"created_at"`     //创建时间
	UpdatedAt     helpers.TimeNormal `json:"updated_at"`     //更新时间
	DeletedAt     gorm.DeletedAt     `json:"deleted_at"`     //删除时间，为 null 则是没删除
}

type Params struct {
	UserId                 int     // 代理id
	Email                  string  // 用户邮箱
	WalletType             int     // 钱包类型：1现货 2合约
	TradingPairId          int     // 交易对id
	TradingPairName        string  // 交易对名称
	TransactionAmount      float64 // 交易金额
	ParentDividend         float64 // 上级获得的分润比例
	WalletStreamType       string  // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	WalletStreamTypeDetail string  // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	Current                int     // 层级默认只处理10层关系
}

// 代理分红

func ParentAgentDividend(params Params) {
	var Parent ParentAgent
	DB := mysql.DB.Debug().Begin()
	DB.Model(models.User{}).
		Where("id", params.UserId). // 上级代理
		Where("is_agent", "1").     // 挑选是代理的
		Where("status", "0").       // 挑选是正常的
		Find(&Parent)
	logger.Info(Parent)

	var AgentDividend float64
	var err error
	if len(Parent.AgentDividend) > 0 {
		AgentDividend, err = strconv.ParseFloat(Parent.AgentDividend, 64) // 代理分红
		if err != nil {
			logger.Error(err)
			return
		}
	} else {
		AgentDividend = 0
	}
	if AgentDividend <= 0 {
		DB.Rollback()
		// 上级代理不存在或者上级红利未设置，不做处理，直接跳过
		logger.Error(errors.New("上级代理红利未设置，不做处理，直接跳过"))
		return
	}
	// 代理所得分红金额
	var Amount float64
	if params.Current == 10 {
		Amount = params.TransactionAmount * AgentDividend * 0.01
	} else {
		Amount = params.TransactionAmount * (params.ParentDividend - AgentDividend) * 0.01
	}
	if Amount <= 0 {
		// 层层瓜分，最终没有钱分了
		logger.Error(errors.New("层层瓜分，最终没有钱分了，不做处理，直接跳过,不继续往下执行"))
		return
	}
	var UsersWallet response.UsersWallet
	DB.Model(models.UsersWallet{}).
		Where("user_id", params.UserId).              // 用户id
		Where("type", params.WalletType).             // 钱包类型：1现货 2合约
		Where("TradingPairId", params.TradingPairId). // 交易对id
		Find(&UsersWallet)
	if UsersWallet.Id <= 0 {
		DB.Rollback()
		logger.Error(errors.New("用户钱包不存在"))
		// 用户钱包不存在
		return
	}
	// 最终可用余额
	Available := UsersWallet.Available + Amount
	DB.Model(models.UsersWallet{}).
		Where("user_id", params.UserId).              // 用户id
		Where("type", params.WalletType).             // 钱包类型：1现货 2合约
		Where("TradingPairId", params.TradingPairId). // 交易对id
		Update("available", Available)                // 更新最终可用余额

	// 记录钱包流水
	var WalletStream models.WalletStream
	WalletStream.TradingPairId = params.TradingPairId       // 交易对ID
	WalletStream.TradingPairName = params.TradingPairName   // 交易对名称
	WalletStream.UserId = params.UserId                     // 用户id
	WalletStream.Email = params.Email                       // 用户邮箱
	WalletStream.Amount = fmt.Sprintf("%v", Amount)         // 流转金额
	WalletStream.HandlingFee = "0"                          // 手续费
	WalletStream.AmountBefore = UsersWallet.Available       // 流转前的余额
	WalletStream.AmountAfter = Available                    // 流转后的余额
	WalletStream.Way = "1"                                  // 流转方式 1 收入 2 支出
	WalletStream.Type = params.WalletStreamType             // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	WalletStream.TypeDetail = params.WalletStreamTypeDetail // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	cErr := DB.Model(WalletStream).Create(&WalletStream).Error
	if cErr != nil {
		DB.Rollback()
		logger.Error(errors.New("添加数据失败" + cErr.Error()))
		return
	}
	// 提交事务
	DB.Commit()
	params.Current-- // 层级递增
	// 只处理10层以内的关系
	if params.Current >= 0 && Parent.ParentId > 0 {
		var data Params
		data.UserId = Parent.ParentId                               // 代理id
		data.Email = Parent.Email                                   // 用户邮箱
		data.WalletType = params.WalletType                         // 钱包类型：1现货 2合约
		data.TradingPairId = params.TradingPairId                   // 交易对id
		data.TradingPairName = params.TradingPairName               // 交易对名称
		data.TransactionAmount = params.TransactionAmount           // 交易金额
		data.ParentDividend = AgentDividend                         // 上级获得的分润比例
		data.WalletStreamType = params.WalletStreamType             // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
		data.WalletStreamTypeDetail = params.WalletStreamTypeDetail // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
		data.Current = params.Current                               // 层级
		ParentAgentDividend(data)
	}
}
