package requests

type ListOptionContractTransaction struct {
	Page       int    `json:"page" form:"page"`
	Limit      int    `json:"limit" form:"limit"`
	OrderBy    string `json:"orderBy" form:"orderBy"`
	Status     string `json:"status" form:"status" validate:"omitempty,min=0"` // 状态：0 交易中 1 已完成 2 已撤单
	CurrencyId string `json:"currency_id" form:"currency_id"`                  // 查询的币种id
}

type OptionContractTransaction struct {
	OptionContractId int     `json:"option_contract_id" form:"option_contract_id" validate:"required"` //期权合约id
	Seconds          int     `json:"seconds" form:"seconds" validate:"required"`                       //秒数
	ProfitRatio      float64 `json:"profit_ratio" form:"profit_ratio" validate:"required"`             //收益率
	Price            string  `json:"price" form:"price" validate:"required"`                           //交易金额
	BuyPrice         string  `json:"buy_price" form:"buy_price" validate:"required"`                   //购买价格
	CurrencyId       string  `json:"currency_id" form:"currency_id" validate:"required"`               // 购买币种
	OrderType        string  `json:"order_type" form:"order_type" validate:"required"`                 // 订单类型：1-买涨 2-买跌
	//ClinchPrice string `json:"clinch_price" form:"clinch_price"` //成交价格
	//CurrencyName string `json:"currency_name" form:"currency_name"` //币种名称 例如：BTC/USDT（币种/交易对）
}
