package models

import (
	"goapi/pkg/config"
	"goapi/pkg/helpers"
	"gorm.io/gorm"
)

// 提币地址

type WalletAddress struct {
	Id        uint64             `json:"id"`
	UserId    int                `json:"user_id"`    //用户id
	Email     string             `json:"email"`      //用户邮箱
	Name      string             `json:"name"`       //地址名称
	Pact      string             `json:"pact"`       //协议： 1-OMNI 2-ERC20 3-TRC20
	Address   string             `json:"address"`    //提币地址
	CreatedAt helpers.TimeNormal `json:"created_at"` //创建时间
	UpdatedAt helpers.TimeNormal `json:"updated_at"` //更新时间
	DeletedAt gorm.DeletedAt     `json:"deleted_at"` //删除时间，为 null 则是没删除
}

func (m *WalletAddress) TableName() string {
	prefix := config.GetString("database.mysql.prefix")
	table := "wallet_address"
	return prefix + table
}
