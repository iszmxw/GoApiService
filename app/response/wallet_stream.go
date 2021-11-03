package response

import (
	"goapi/pkg/helpers"
)

// WalletStream 钱包流水
type WalletStream struct {
	Id              int                `json:"id"`
	TradingPairId   string             `json:"trading_pair_id"`   //交易对ID
	TradingPairName string             `json:"trading_pair_name"` //交易对名称
	UserId          string             `json:"user_id"`           //用户id
	Email           string             `json:"email"`             //用户邮箱
	Amount          string             `json:"amount"`            //流转金额
	AmountBefore    string             `json:"amount_before"`     //流转前的余额
	AmountAfter     string             `json:"amount_after"`      //流转后的余额
	Way             string             `json:"way"`               //流转方式 1 收入 2 支出
	Type            string             `json:"type"`              //流转类型 1 币币交易 2 永续合约 3 期权合约 4 申购交易 5 划转 6 充值 7 提现 8 冻结
	TypeDetail      string             `json:"type_detail"`       //流转详细类型  1 USDT充值  2 银行卡充值  3 币币交易手续费  4 永续合约手续费  5 期权合约手续费  6 币币账户划转到合约账户  7 合约账户划转到币币账户  8 申购冻结  9 币币交易  10 永续合约  11 期权合约
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       helpers.TimeNormal `json:"deleted_at"`        //删除时间，为 null 则是没删除
}
