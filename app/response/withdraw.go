package response

import (
	"goapi/pkg/helpers"
)

// Withdraw 用户提现订单记录表
type Withdraw struct {
	Id              uint64             `json:"id"`
	Email           string             `json:"email"`            //邮箱
	Address         string             `json:"address"`          //提币地址
	CurrencyId      int8               `json:"currency_id"`      //提现币种id
	CurrencyName    string             `json:"currency_name"`    //提现币种名称
	CoinChainType   int8               `json:"coin_chain_type"`  //币链类型 1-OMNI 2-ERC20 3-TRC20
	WithdrawNum     string             `json:"withdraw_num"`     //提现数量
	HandlingFee     string             `json:"handling_fee"`     //手续费
	ActuallyArrived string             `json:"actually_arrived"` //实际到账
	Status          string             `json:"status"`           //状态：0-未确认：1-已确认 2-已拒绝
	Remark          string             `json:"remark"`           //备注
	CreatedAt       helpers.TimeNormal `json:"created_at"`       //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`       //更新时间
	DeletedAt       helpers.TimeNormal `json:"deleted_at"`       //删除时间，为 null 则是没删除
}
