package requests

type ListRecharge struct {
	Page            int    `json:"page" form:"page"`
	Limit           int    `json:"limit" form:"limit"`
	OrderBy         string `json:"orderBy" form:"orderBy"`
	Status          string `json:"status" form:"status" validate:"omitempty,min=0"` // 状态：0 交易中 1 已完成 2 已撤单
	TradingPairName string `json:"trading_pair_name" form:"trading_pair_name"`      // 查询的币种id
}

// Recharge 用户充值订单记录表
type Recharge struct {
	Type            string `json:"type" form:"type" validate:"required,numeric"`                       // 类型：1-OMNI  2-ERC20  3-TRC20
	Address         string `json:"address" form:"address" validate:"required"`                         // 充币地址
	RechargeNum     string `json:"recharge_num" form:"recharge_num" validate:"required,numeric"`       // 充值数量
	AccountNo       string `json:"account_no" form:"account_no" validate:"required"`                   // 银行卡号
	BankCode        string `json:"bank_code" form:"bank_code" validate:"required"`                     // 银行编码
	TradingPairId   string `json:"trading_pair_id" form:"trading_pair_id" validate:"required,numeric"` // 充值的交易对id
	TradingPairName string `json:"trading_pair_name" form:"trading_pair_name" validate:"required"`     // 充值的交易对name
	TradingPairType string `json:"trading_pair_type" form:"trading_pair_type" validate:"required"`     // 充值的交易对类型
}
