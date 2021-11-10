package response

import (
	"goapi/pkg/helpers"
)

type Currency struct {
	Id                   uint               `json:"id"`                     // 主键id
	Name                 string             `json:"name"`                   // 名称
	TradingPairId        int                `json:"trading_pair_id"`        // 交易对
	TradingPairName      string             `json:"trading_pair_name"`      // 交易对名称
	KLineCode            string             `json:"k_line_code"`            // K线图代码
	DecimalScale         int                `json:"decimal_scale"`          // 自有币位数
	Type                 string             `json:"type"`                   // 交易显示：（1币币交易，2永续合约，3期权合约）
	Sort                 int8               `json:"sort"`                   // 排序
	IsHidden             int8               `json:"is_hidden"`              // 是否展示：0-否，1-展示
	FeePerpetualContract float64            `json:"fee_perpetual_contract"` // 永续合约手续费%
	FeeOptionContract    float64            `json:"fee_option_contract"`    // 期权合约手续费%
	CreatedAt            helpers.TimeNormal `json:"created_at"`             // 创建时间
	UpdatedAt            helpers.TimeNormal `json:"updated_at"`             // 更新时间
	DeletedAt            helpers.TimeNormal `json:"deleted_at"`             // 删除时间，为 null 则是没删除
}
