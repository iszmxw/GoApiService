package requests

// TradingPair 交易对id
type TradingPair struct {
	TradingPairId string `json:"trading_pair_id" form:"trading_pair_id" validate:"required,numeric"` //提现的交易对id
}

// ListAssetsStream 查询资产流水
type ListAssetsStream struct {
	TradingPairId string `json:"trading_pair_id" form:"trading_pair_id" validate:"omitempty,numeric"` //交易对ID
	OrderType     string `json:"order_type" form:"order_type" validate:"omitempty,numeric"`           //订单类型 1 币币交易 2 永续合约 3 期权合约
	Time          string `json:"time" form:"time" validate:"omitempty,numeric"`                       //订单时间最近七天（7, 15, 30）
}
