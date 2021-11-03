package models

import (
	"goapi/pkg/helpers"
	"gorm.io/gorm"
)

// ApplyBuySetup 申购币种设置
type ApplyBuySetup struct {
	Id              int                `json:"id"`                //主键id
	Name            string             `json:"name"`              //币种名称
	TradingPairId   int                `json:"trading_pair_id"`   //交易对ID
	TradingPairName string             `json:"trading_pair_name"` //交易对名称
	Ratio           float64            `json:"ratio"`             //购买比例
	EstimatedTime   helpers.TimeNormal `json:"estimated_time"`    //预计上线时间
	StartTime       helpers.TimeNormal `json:"start_time"`        //开始申购时间
	EndTime         helpers.TimeNormal `json:"end_time"`          //结束申购时间
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       gorm.DeletedAt     `json:"deleted_at"`        //删除时间，为 null 则是没删除
}
