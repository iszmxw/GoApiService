package models

import (
	"fmt"
	"goapi/app/response"
	"goapi/pkg/helpers"
	"gorm.io/gorm"
)

// AssetsStream 资产流水记录表
type AssetsStream struct {
	Id              int                `json:"id"`
	UserId          int                `json:"user_id"`           //用户id
	Email           string             `json:"email"`             //用户邮箱
	CurrencyId      int                `json:"currency_id"`       //交易币种
	CurrencyName    string             `json:"currency_name"`     //币种名称 例如：BTC/USDT（币种/交易对）
	TradingPairId   int                `json:"trading_pair_id"`   //交易对ID
	TradingPairName string             `json:"trading_pair_name"` //交易对名称
	OrderType       int                `json:"order_type"`        //订单类型 1 币币交易 2 永续合约 3 期权合约
	OrderId         int                `json:"order_id"`          //订单id
	OrderNum        string             `json:"order_num"`         //订单编号
	OrderTime       helpers.TimeNormal `json:"order_time"`        //订单时间
	OrderPrice      string             `json:"order_price"`       //交易金额
	Status          string             `json:"status"`            //状态：0 交易中 1 已完成 2 已撤单
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       gorm.DeletedAt     `json:"deleted_at"`        //删除时间，为 null 则是没删除
}

// SetAddData 设置添加数据
func (m *AssetsStream) SetAddData(OrderType int, addData interface{}, Currency response.Currency) *AssetsStream {
	data := m
	// 币币交易数据
	if OrderType == 1 {
		data.OrderType = 1
		data.UserId = addData.(CurrencyTransaction).UserId
		data.Email = addData.(CurrencyTransaction).Email
		data.CurrencyId = addData.(CurrencyTransaction).CurrencyId
		data.CurrencyName = addData.(CurrencyTransaction).CurrencyName
		data.TradingPairId = Currency.TradingPairId
		data.TradingPairName = Currency.TradingPairName
		data.OrderId = addData.(CurrencyTransaction).Id
		data.OrderNum = addData.(CurrencyTransaction).OrderNumber
		data.OrderTime = addData.(CurrencyTransaction).CreatedAt
		data.OrderPrice = fmt.Sprintf("%v", addData.(CurrencyTransaction).OrderPrice)
		data.Status = addData.(CurrencyTransaction).Status
	}

	// 永续合约交易数据
	if OrderType == 2 {
		data.OrderType = 2
		data.UserId = addData.(PerpetualContractTransaction).UserId
		data.Email = addData.(PerpetualContractTransaction).Email
		data.CurrencyId = addData.(PerpetualContractTransaction).CurrencyId
		data.CurrencyName = addData.(PerpetualContractTransaction).CurrencyName
		data.TradingPairId = Currency.TradingPairId
		data.TradingPairName = Currency.TradingPairName
		data.OrderId = addData.(PerpetualContractTransaction).Id
		data.OrderNum = addData.(PerpetualContractTransaction).OrderNumber
		data.OrderTime = addData.(PerpetualContractTransaction).CreatedAt
		data.OrderPrice = addData.(PerpetualContractTransaction).Price // 期权合约
		data.Status = addData.(PerpetualContractTransaction).Status

	}

	// 期权合约交易数据
	if OrderType == 3 {
		data.OrderType = 3
		data.UserId = addData.(OptionContractTransaction).UserId
		data.Email = addData.(OptionContractTransaction).Email
		data.CurrencyId = addData.(OptionContractTransaction).CurrencyId
		data.CurrencyName = addData.(OptionContractTransaction).CurrencyName
		data.TradingPairId = Currency.TradingPairId
		data.TradingPairName = Currency.TradingPairName
		data.OrderId = addData.(OptionContractTransaction).Id
		data.OrderNum = addData.(OptionContractTransaction).OrderNumber
		data.OrderTime = addData.(OptionContractTransaction).CreatedAt
		data.OrderPrice = addData.(OptionContractTransaction).Price
		data.Status = addData.(OptionContractTransaction).Status
	}

	return data
}
