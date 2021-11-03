package response

import (
	"goapi/pkg/helpers"
)

type OptionContractTransaction struct {
	Id               uint               `json:"id"`                 //主键id
	UserId           int                `json:"user_id"`            //用户id
	Email            string             `json:"email"`              //用户邮箱
	OrderNumber      string             `json:"order_number"`       //订单号
	OptionContractId int                `json:"option_contract_id"` //期权合约id
	Seconds          int64              `json:"seconds"`            //秒数
	Time             int64              `json:"time"`               //秒数
	ProfitRatio      float64            `json:"profit_ratio"`       //收益率
	Price            string             `json:"price"`              //交易金额
	BuyPrice         string             `json:"buy_price"`          //购买价格
	HandleFee        string             `json:"handle_fee"`         //期权合约手续费
	ClinchPrice      string             `json:"clinch_price"`       //成交价格
	CurrencyId       int                `json:"currency_id"`        //购买币种id
	CurrencyName     string             `json:"currency_name"`      //币种名称 例如：BTC/USDT（币种/交易对）
	TradingPairId    string             `json:"trading_pair_id"`    //交易对id
	TradingPairName  string             `json:"trading_pair_name"`  //交易对名称
	OrderType        string             `json:"order_type"`         //订单类型：1-买涨 2-买跌
	OrderResult      string             `json:"order_result"`       //交割结果：1-涨 2-跌
	ResultProfit     string             `json:"result_profit"`      //预计最终盈利金额
	Status           int                `json:"status"`             //状态：0 交易中 1 已完成 2 已撤单
	Remark           string             `json:"remark"`             //备注
	CreatedAt        helpers.TimeNormal `json:"created_at"`         //创建时间
	UpdatedAt        helpers.TimeNormal `json:"updated_at"`         //更新时间
	DeletedAt        helpers.TimeNormal `json:"deleted_at"`         //删除时间，为 null 则是没删除
}

// 通过订单id关联查询出K线图信息

type OptionTransactionKline struct {
	Id               uint               `json:"id"`                 //主键id
	UserId           int                `json:"user_id"`            //用户id
	Email            string             `json:"email"`              //用户邮箱
	OrderNumber      string             `json:"order_number"`       //订单号
	OptionContractId int                `json:"option_contract_id"` //期权合约id
	Seconds          int                `json:"seconds"`            //秒数
	ProfitRatio      float64            `json:"profit_ratio"`       //收益率
	Price            float64            `json:"price"`              //交易金额
	BuyPrice         string             `json:"buy_price"`          //购买价格
	HandleFee        float64            `json:"handle_fee"`         //期权合约手续费
	ClinchPrice      string             `json:"clinch_price"`       //成交价格
	CurrencyId       int                `json:"currency_id"`        //购买币种id
	CurrencyName     string             `json:"currency_name"`      //币种名称 例如：BTC/USDT（币种/交易对）
	TradingPairId    string             `json:"trading_pair_id"`    //交易对id
	TradingPairName  string             `json:"trading_pair_name"`  //交易对名称
	KLineCode        string             `json:"k_line_code"`        //k线图
	OrderType        string             `json:"order_type"`         //订单类型：1-买涨 2-买跌
	OrderResult      string             `json:"order_result"`       //交割结果：1-涨 2-跌
	Status           string             `json:"status"`             //状态：0 交易中 1 已完成 2 已撤单
	Remark           string             `json:"remark"`             //备注
	CreatedAt        helpers.TimeNormal `json:"created_at"`         //创建时间
	UpdatedAt        helpers.TimeNormal `json:"updated_at"`         //更新时间
	DeletedAt        helpers.TimeNormal `json:"deleted_at"`         //删除时间，为 null 则是没删除
}
