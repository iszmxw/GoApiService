package requests

// TradingPair 交易对id
type TradingPair struct {
	Type            string `json:"type" form:"type" validate:"required,numeric"`                        // 1现货 2合约
	TradingPairId   string `json:"trading_pair_id" form:"trading_pair_id" validate:"omitempty,numeric"` //提现的交易对id
	TradingPairName string `json:"trading_pair_name" form:"trading_pair_name"`                          //提现的交易对name
}

// ListAssetsStream 查询资产流水
type ListAssetsStream struct {
	TradingPairId string `json:"trading_pair_id" form:"trading_pair_id" validate:"omitempty,numeric"` // 交易对ID
	OrderType     string `json:"order_type" form:"order_type" validate:"omitempty,numeric"`           // 流转类型 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	StartTime     string `json:"start_time" form:"start_time"`                                        // 开始时间
	EndTime       string `json:"end_time" form:"end_time"`                                            // 结束时间
	Page          int    `json:"page" form:"page" validate:"omitempty,numeric"`                       // 第几页
	Limit         int    `json:"limit" form:"limit" validate:"omitempty,numeric"`                     // 获取多少条
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

//  AssetsStream 查询个人资产

type AssetsStream struct {
	Type string `json:"type" form:"type" validate:"omitempty,numeric"` // 1现货 2合约
}
