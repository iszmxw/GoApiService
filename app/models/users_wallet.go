package models

import (
	"gorm.io/gorm"
	"time"
)

type UsersWallet struct {
	Id              uint           // 主键id
	UserId          int            // 用户id
	TradingPairId   int            //交易对id
	TradingPairName string         //交易对名称
	Address         string         // 钱包地址
	Status          int            // 0正常 1锁定
	Available       float64        // 可用
	Freeze          float64        // 冻结
	TotalAssets     float64        // 总资产
	CreatedAt       time.Time      // 创建时间
	UpdatedAt       time.Time      // 更新时间
	DeletedAt       gorm.DeletedAt // 删除时间，为 null 则是没删除
}
