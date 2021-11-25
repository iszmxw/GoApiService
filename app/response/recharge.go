package response

import (
	"goapi/pkg/helpers"
)

// Recharge 用户充值订单记录表
type Recharge struct {
	Id              int                `json:"id"`
	TopUpType       int                `json:"top_up_type"`       //充值类型：1-USDT 2-银行卡
	PayId           int                `json:"pay_id"`            //支付订单id
	UserId          int                `json:"user_id"`           //用户id
	Email           string             `json:"email"`             //用户邮箱
	Address         string             `json:"address"`           //充币地址
	Type            int8               `json:"type"`              //类型：1-OMNI  2-ERC20  3-TRC20
	TradingPairId   int8               `json:"trading_pair_id"`   //充值的交易对id
	TradingPairName string             `json:"trading_pair_name"` //充值的交易对name
	TradingPairType string             `json:"trading_pair_type"` //充值的交易对类型 1现货 2合约
	RechargeNum     string             `json:"recharge_num"`      //充值数量
	Status          string             `json:"status"`            //状态：0-未确认：1-已确认
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       helpers.TimeNormal `json:"deleted_at"`        //删除时间，为 null 则是没删除
}
