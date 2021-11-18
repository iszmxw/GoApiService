package models

import (
	"gorm.io/gorm"
	"time"
)

// option_contract  秒合约配置

type OptionContract struct {
	Id          uint           `gorm:"column:id"`           // 主键id
	Seconds     int            `gorm:"column:seconds"`      // 秒数
	Status      int8           `gorm:"column:status"`       // 状态:0.禁用,1.启用
	ProfitRatio float64        `gorm:"column:profit_ratio"` // 收益率
	Minimum     float64        `gorm:"column:minimum"`      // 最低入场
	CreatedAt   time.Time      `gorm:"column:created_at"`   // 创建时间
	UpdatedAt   time.Time      `gorm:"column:updated_at"`   // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`   // 删除时间，为 null 则是没删除
}
