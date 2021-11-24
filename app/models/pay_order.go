package models

import (
	"goapi/pkg/config"
	"goapi/pkg/helpers"
)

// PayOrder 支付订单记录
type PayOrder struct {
	Id            uint64             `json:"id"`
	UserId        int                `json:"user_id"`         //用户id
	Email         string             `json:"email"`           //用户邮箱
	AccountNo     string             `json:"account_no"`      //用户银行卡号
	BankCode      string             `json:"bank_code"`       //银行卡代码
	Amount        float64            `json:"amount"`          //订单金额金额
	PayAmount     float64            `json:"pay_amount"`      //支付金额
	Fee           float64            `json:"fee"`             //手续费
	OrderNum      string             `json:"order_num"`       //上送订单号
	SystemOrder   string             `json:"system_order"`    //平台订单号
	Status        int8               `json:"status"`          //订单状态，0：待支付，1：支付成功 4:支付失败
	SuccessTime   helpers.TimeNormal `json:"success_time"`    //成功支付时间
	Payurl        string             `json:"payurl"`          //支付跳转地址
	Qrcode        string             `json:"qrcode"`          //二维码信息
	PayeeNo       string             `json:"payee_no"`        //付款账户
	PayeeName     string             `json:"payee_name"`      //付款账户名称
	PayeeBankCode string             `json:"payee_bank_code"` //付款账户银行代码
	CreatedAt     helpers.TimeNormal `json:"created_at"`      //创建时间
	UpdatedAt     helpers.TimeNormal `json:"updated_at"`      //更新时间
	DeletedAt     helpers.TimeNormal `json:"deleted_at"`      //删除时间，为 null 则是没删除
}

func (m *PayOrder) TableName() string {
	prefix := config.GetString("database.mysql.prefix")
	table := "pay_order"
	return prefix + table
}
