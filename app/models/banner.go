package models

import (
	"goapi/pkg/helpers"
)

// banner管理

type Banner struct {
	Id        uint               `json:"id"`         // 主键id
	Name      string             `json:"name"`       // 广告名称
	Images    string             `json:"images"`     // 广告图片
	Redirect  string             `json:"redirect"`   // 跳转地址
	Sort      int                `json:"sort"`       // 排序
	Type      string             `json:"type"`       // 广告位置：1-pc 2-h5 3-app
	CreatedAt helpers.TimeNormal `json:"created_at"` // 创建时间
	UpdatedAt helpers.TimeNormal `json:"updated_at"` // 更新时间
	DeletedAt helpers.TimeNormal `json:"deleted_at"` // 删除时间，为 null 则是没删除
}
