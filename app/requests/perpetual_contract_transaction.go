package requests

type ListPerpetualContractTransaction struct {
	Page       int    `json:"page" form:"page"`
	Limit      int    `json:"limit" form:"limit"`
	OrderBy    string `json:"orderBy" form:"orderBy"`
	Status     string `json:"status" form:"status" validate:"omitempty,min=0"` // 状态：0 交易中 1 已完成 2 已撤单
	CurrencyId string `json:"currency_id" form:"currency_id"`                  // 查询的币种id
}

type PerpetualContractTransaction struct {
	CurrencyId      string `json:"currency_id" form:"currency_id" validate:"required,numeric"`   // 币种
	OrderType       string `json:"order_type" form:"order_type" validate:"required"`             // 订单类型：1-限价 2-市价
	LimitPrice      string `json:"limit_price" form:"limit_price"`                               // 当前限价
	TransactionType string `json:"transaction_type" form:"transaction_type" validate:"required"` // 交易类型：1-开多 2-开空
	EntrustNum      string `json:"entrust_num" form:"entrust_num" validate:"required"`           // 委托量
	EntrustPrice    string `json:"entrust_price" form:"entrust_price" validate:"required"`       // 委托价格
	EnsureAmount    string `json:"ensure_amount" form:"ensure_amount" validate:"required"`       // 保证金
	HandNum         string `json:"hand_num" form:"hand_num" validate:"required"`                 // 手数值
	Multiple        string `json:"multiple" form:"multiple" validate:"required"`                 // 倍数值
}
