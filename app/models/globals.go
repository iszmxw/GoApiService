package models

import (
	"gorm.io/gorm"
	"time"
)

// Globals 全局配置管理
type Globals struct {
	Id        uint32         `json:"id"`         //主键id
	Fields    string         `json:"fields"`     //字段名称
	Value     string         `json:"value"`      //字段值
	CreatedAt time.Time      `json:"created_at"` //创建时间
	UpdatedAt time.Time      `json:"updated_at"` //更新时间
	DeletedAt gorm.DeletedAt `json:"deleted_at"` //删除时间，为 null 则是没删除
}

type WithdrawalFees struct {
	Id     uint32  `json:"id"`     //主键id
	Fields string  `json:"fields"` //字段名称
	Value  float64 `json:"value"`  //字段值
}

type MinAmount struct {
	Id     uint32  `json:"id"`     //主键id
	Fields string  `json:"fields"` //字段名称
	Value  float64 `json:"value"`  //字段值
}
