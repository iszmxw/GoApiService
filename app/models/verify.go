package models

import (
	"goapi/pkg/config"
	"gorm.io/gorm"
	"time"
)

type Verify struct {
	UserId        int            `json:"user_id" form:"user_id"`                 //用户id从token中拿
	IdentityCard  string         `json:"identity_card" form:"identity_card"`     //身份证信息
	Status        int            `json:"status" form:"status"`                   //状态0未提交 1初级验证  2高级验证
	Email         string         `json:"email " form:"email"`                    //邮箱
	FullName      string         `json:"full_name" form:"full_name"`             //用户真实姓名
	ImgCardBehind string         `json:"img_card_behind" form:"img_card_behind"` //身份证后面图片
	ImgCardFront  string         `json:"img_card_front"  form:"img_card_front"`  //身份证前面图片
	ImgBankFront  string         `json:"img_bank_front" form:"img_bank_front"`   //银行卡前面
	ImgBankBehind string         `json:"img_bank_behind" form:"img_bank_behind"` //银行卡后面
	CreatedAt     time.Time      `gorm:"column:created_at"`                      // 创建时间
	UpdatedAt     time.Time      `gorm:"column:updated_at"`                      // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`                      // 删除时间，为 null 则是没删除
}

func (m *Verify) TableName() string {
	prefix := config.GetString("database.mysql.prefix")
	table := "user_img"
	return prefix + table
}
