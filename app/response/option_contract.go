package response

import (
	"goapi/pkg/helpers"
)

// option_contract  秒合约配置

type OptionContract struct {
	Id          int                `json:"id"`           // 主键id
	Seconds     int                `json:"seconds"`      // 秒数
	Status      int                `json:"status"`       // 状态:0.禁用,1.启用
	ProfitRatio float64            `json:"profit_ratio"` // 收益率
	CreatedAt   helpers.TimeNormal `json:"created_at"`   // 创建时间
	UpdatedAt   helpers.TimeNormal `json:"updated_at"`   // 更新时间
	DeletedAt   helpers.TimeNormal `json:"deleted_at"`   // 删除时间，为 null 则是没删除
}
