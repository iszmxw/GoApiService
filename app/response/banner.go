package response

import (
	"goapi/pkg/helpers"
)

// banner管理

type Banner struct {
	Id        uint               `json:"id" form:"id"`                 //主键id
	Name      string             `json:"name" form:"name"`             //广告名称
	Images    string             `json:"images" form:"images"`         //广告图片
	Redirect  string             `json:"redirect" form:"redirect"`     //跳转地址
	Sort      int                `json:"sort" form:"sort"`             //排序
	Type      string             `json:"type" form:"type"`             //广告位置：1-pc 2-h5 3-app
	CreatedAt helpers.TimeNormal `json:"created_at" form:"created_at"` //创建时间
	UpdatedAt helpers.TimeNormal `json:"updated_at" form:"updated_at"` //更新时间
	DeletedAt helpers.TimeNormal `json:"deleted_at" form:"deleted_at"` //删除时间，为 null 则是没删除
}
