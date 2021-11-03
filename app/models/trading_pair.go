package models

import (
	"gorm.io/gorm"
	"time"
)

type TradingPair struct {
	Id        uint           `gorm:"column:id" json:"id"`                 //主键id
	Name      string         `gorm:"column:name" json:"name"`             //交易对名称
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"` //创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"` //更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"` //删除时间，为 null 则是没删除
}
