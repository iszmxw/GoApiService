package response

import (
	"goapi/pkg/helpers"
)

type CurrencyTransaction struct {
	Id              uint               `json:"id"`               //主键id
	UserId          int                `json:"user_id"`          //用户id
	Email           string             `json:"email"`            //用户邮箱
	OrderNumber     string             `json:"order_number"`     //订单号
	CurrencyId      int                `json:"currency_id"`      //币种
	Name            string             `json:"name"`             //币种名称
	TradingPairId   int                `json:"trading_pair_id"`  //交易对id
	CurrencyName    string             `json:"currency_name"`    //币种名称 例如：BTC/USDT（币种/交易对）
	EntrustNum      string             `json:"entrust_num"`      // 委托量
	OrderType       string             `json:"order_type"`       //挂单类型：1-限价 2-市价
	LimitPrice      string             `json:"limit_price"`      //当前限价
	TransactionType string             `json:"transaction_type"` //订单方向：1-买入 2-卖出
	OrderPrice      float64            `json:"price"`            //订单额
	Status          string             `json:"status"`           //状态：0 交易中 1 已完成 2 已撤单
	Remark          string             `json:"remark"`           //备注
	CreatedAt       helpers.TimeNormal `json:"created_at"`       //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`       //更新时间
	DeletedAt       helpers.TimeNormal `json:"deleted_at"`       //删除时间，为 null 则是没删除
}
