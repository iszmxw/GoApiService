package requests

// AddWithdraw 添加数据
type AddWithdraw struct {
	AddressId       string `json:"address_id" form:"address_id" validate:"required"`                     // 提币地址 id
	Type            string `json:"type" form:"type" validate:"required,numeric"`                         // 币链类型 1-OMNI 2-ERC20 3-TRC20
	WithdrawNum     string `json:"withdraw_num" form:"withdraw_num" validate:"required,numeric"`         // 提现数量
	ActuallyArrived string `json:"actually_arrived" form:"actually_arrived" validate:"required,numeric"` // 实际到账
}

// ListWithdraw 查询数据
type ListWithdraw struct {
	Page            int    `json:"page" form:"page"`
	Limit           int    `json:"limit" form:"limit"`
	OrderBy         string `json:"orderBy" form:"orderBy"`
	Status          string `json:"status" form:"status" validate:"omitempty,min=0"` // 状态：0 交易中 1 已完成 2 已撤单
	TradingPairName string `json:"trading_pair_name" form:"trading_pair_name"`      // 交易对名称
}
