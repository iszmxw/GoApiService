package response

import (
	"goapi/pkg/helpers"
)

// 永续合约，合约信息

type PerpetualContract struct {
	Id         int                `json:"id"`          // 主键id
	CurrencyId int                `json:"currency_id"` // 交易对id（从币种管理获取）
	Multiple   string             `json:"multiple"`    // 倍数：10、25、50、100
	Bail       string             `json:"bail"`        // 保证金占比：100、50、25、10
	Ratio      string             `json:"ratio"`       // 张数比例：1：？USDT
	CreatedAt  helpers.TimeNormal `json:"created_at"`  // 创建时间
	UpdatedAt  helpers.TimeNormal `json:"updated_at"`  // 更新时间
	DeletedAt  helpers.TimeNormal `json:"deleted_at"`  // 删除时间，为 null 则是没删除
}
