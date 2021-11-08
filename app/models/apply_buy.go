package models

import (
	"goapi/pkg/helpers"
	"gorm.io/gorm"
)

// ApplyBuy 申购币种记录
type ApplyBuy struct {
	Id              int                `json:"id"`                // 主键id
	UserId          string             `json:"user_id"`           // 用户id
	Email           string             `json:"email"`             // 申购用户
	GetCurrencyId   int                `json:"get_currency_id"`   // 购买币种id
	GetCurrencyName string             `json:"get_currency_name"` // 购买币种名称
	GetCurrencyNum  float64            `json:"get_currency_num"`  // 购买数量
	TradingPairId   int                `json:"trading_pair_id"`   // 交易对ID
	TradingPairName string             `json:"trading_pair_name"` // 交易对名称
	IssuePrice      float64            `json:"issue_price"`       // 发行价 1 = 多少个USDT
	CreatedAt       helpers.TimeNormal `json:"created_at"`        // 创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        // 更新时间
	DeletedAt       gorm.DeletedAt     `json:"deleted_at"`        // 删除时间，为 null 则是没删除
}
