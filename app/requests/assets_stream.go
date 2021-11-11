package requests

// TradingPair 交易对id
type TradingPair struct {
	Type          string `json:"type" form:"type" validate:"required,numeric"`                       // 1现货 2合约
	TradingPairId string `json:"trading_pair_id" form:"trading_pair_id" validate:"required,numeric"` //提现的交易对id
}

// ListAssetsStream 查询资产流水
type ListAssetsStream struct {
	TradingPairId string `json:"trading_pair_id" form:"trading_pair_id" validate:"omitempty,numeric"` // 交易对ID
	OrderType     string `json:"order_type" form:"order_type" validate:"omitempty,numeric"`           // 流转类型 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	Time          string `json:"time" form:"time" validate:"omitempty,numeric"`                       // 订单时间最近七天（7, 15, 30）
}

// WalletAddressAdd 提币地址配置
type WalletAddressAdd struct {
	Name    string `json:"name" form:"name" validate:"required"`       // 地址名称
	Pact    string `json:"pact" form:"pact" validate:"required"`       // 协议： 1-OMNI 2-ERC20 3-TRC20
	Address string `json:"address" form:"address" validate:"required"` // 提币地址
}

// WalletAddressDel 提币地址删除
type WalletAddressDel struct {
	Id string `json:"id" form:"id" validate:"required"` // id
}
