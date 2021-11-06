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
	HandlingFee     string             `json:"handling_fee"`      //手续费
	AmountBefore    float64            `json:"amount_before"`     //流转前的余额
	AmountAfter     float64            `json:"amount_after"`      //流转后的余额
	Way             string             `json:"way"`               //流转方式 1 收入 2 支出
	Type            string             `json:"type"`              //流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	TypeDetail      string             `json:"type_detail"`       //流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
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
	data.Type = Type             // 流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	data.TypeDetail = TypeDetail // 流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
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
