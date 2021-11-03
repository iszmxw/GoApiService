package models

import (
	"errors"
	"goapi/app/response"
	"goapi/pkg/config"
	"goapi/pkg/helpers"
	"gorm.io/gorm"
	"strconv"
)

// WalletStream 钱包流水
type WalletStream struct {
	Id              uint64             `json:"id"`
	TradingPairId   int                `json:"trading_pair_id"`   //交易对ID
	TradingPairName string             `json:"trading_pair_name"` //交易对名称
	UserId          int                `json:"user_id"`           //用户id
	Email           string             `json:"email"`             //用户邮箱
	Amount          string             `json:"amount"`            //流转金额
	AmountBefore    float64            `json:"amount_before"`     //流转前的余额
	AmountAfter     float64            `json:"amount_after"`      //流转后的余额
	Way             string             `json:"way"`               //流转方式 1 收入 2 支出
	Type            string             `json:"type"`              //流转类型 1 币币交易 2 永续合约 3 期权合约 4 申购交易 5 划转 6 充值 7 提现 8 冻结
	TypeDetail      string             `json:"type_detail"`       //流转详细类型  1 USDT充值  2 银行卡充值  3 币币交易手续费  4 永续合约手续费  5 期权合约手续费  6 币币账户划转到合约账户  7 合约账户划转到币币账户  8 申购冻结  9 币币交易  10 永续合约  11 期权合约
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       gorm.DeletedAt     `json:"deleted_at"`        //删除时间，为 null 则是没删除
}

func (m *WalletStream) TableName() string {
	prefix := config.GetString("database.mysql.prefix")
	table := "wallet_stream"
	return prefix + table
}

// SetAddData 生成钱包流水数据结构
func (m *WalletStream) SetAddData(Way, Type, TypeDetail string, addData interface{}, Currency response.Currency, UsersWallet response.UsersWallet) (*WalletStream, error) {
	data := m
	data.Way = Way               // 流转方式 1 收入 2 支出
	data.Type = Type             // 流转类型 1 币币交易 2 永续合约 3 期权合约 4 申购交易 5 划转 6 充值 7 提现 8 冻结
	data.TypeDetail = TypeDetail // 流转详细类型  1 USDT充值  2 银行卡充值  3 币币交易手续费  4 永续合约手续费  5 期权合约手续费  6 币币账户划转到合约账户  7 合约账户划转到币币账户  8 申购冻结  9 币币交易  10 永续合约  11 期权合约
	switch data.Type {
	case "1":
		// 1-币币交易
		{
			data.UserId = addData.(CurrencyTransaction).UserId
			data.Email = addData.(CurrencyTransaction).Email
			data.TradingPairId = Currency.TradingPairId
			data.TradingPairName = Currency.TradingPairName
			data.Amount = addData.(CurrencyTransaction).Price // 流转金额
			data.AmountBefore = UsersWallet.Available         // 流转前的余额
			Amount, err := strconv.ParseFloat(data.Amount, 64)
			if err != nil {
				return nil, err
			}
			data.AmountAfter = Amount + UsersWallet.Available // 流转后的余额
		}
		break
	case "2":
		// 2-永续合约
		{

			data.UserId = addData.(PerpetualContractTransaction).UserId
			data.Email = addData.(PerpetualContractTransaction).Email
			data.TradingPairId = Currency.TradingPairId
			data.TradingPairName = Currency.TradingPairName
			data.Amount = addData.(PerpetualContractTransaction).Price // 流转金额
			data.AmountBefore = UsersWallet.Available                  // 流转前的余额
			Amount, err := strconv.ParseFloat(data.Amount, 64)
			if err != nil {
				return nil, err
			}
			data.AmountAfter = Amount + UsersWallet.Available // 流转后的余额
		}
		break
	case "3":
		// 3-期权合约
		{
			data.UserId = addData.(OptionContractTransaction).UserId
			data.Email = addData.(OptionContractTransaction).Email
			data.TradingPairId = Currency.TradingPairId
			data.TradingPairName = Currency.TradingPairName
			data.Amount = addData.(OptionContractTransaction).Price // 流转金额
			data.AmountBefore = UsersWallet.Available               // 流转前的余额
			Amount, err := strconv.ParseFloat(data.Amount, 64)
			if err != nil {
				return nil, err
			}
			data.AmountAfter = Amount + UsersWallet.Available // 流转后的余额
		}
		break
	default:
		return nil, errors.New("流转类型值异常")
	}
	return data, nil
}
