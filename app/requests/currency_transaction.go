package requests

type CurrencyTransaction struct {
	CurrencyId      string `json:"currency_id" form:"currency_id" validate:"required"`           // 币种
	//EntrustNum      string `json:"entrust_num" form:"entrust_num" validate:"required"`           // todo::收集无用--委托量，委托（20%，50%，75%，100%）
	LimitPrice      string `json:"limit_price" form:"limit_price"`                               // 当前限价
	ClinchNum       string `json:"clinch_num" form:"clinch_num" validate:"required"`             // 成交量
	TransactionType string `json:"transaction_type" form:"transaction_type" validate:"required"` // 订单方向：1-买入 2-卖出
	OrderType       string `json:"order_type" form:"order_type" validate:"required,gte=1,lte=2"` // 挂单类型：1-限价 2-市价
}

// CancelOrder 撤单
type CancelOrder struct {
	Id string `json:"id" form:"id" validate:"required,numeric"` // id
}

// Liquidation 平仓
type Liquidation struct {
	Id          string  `json:"id" form:"id" validate:"required,numeric"`                   // id
	Liquidation float64 `json:"liquidation" form:"liquidation" validate:"required,numeric"` // 平仓
}
