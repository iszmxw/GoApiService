package requests

type Transfer struct {
	Type          string `json:"type" form:"type" validate:"required"`                       // 1 从现货账户划转到合约账户  2 从合约账户划转到现货账户
	TradingPairId string `json:"trading_pair_id" form:"trading_pair_id" validate:"required"` // 币种 id
	Num           string `json:"num" form:"num" validate:"required,numeric"`                 // 划转数量
}
