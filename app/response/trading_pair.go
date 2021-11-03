package response

type TradingPair struct {
	Id   int    `json:"id"`   // 主键id
	Name string `json:"name"` // 交易对名称
}
