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
	Type            string             `json:"type"`              //流转类型 0 未知 1 充值 2 提现 3 划转 4 快捷买币 5 空投 6 现货 7 合约 8 期权 9 手续费
	TypeDetail      string             `json:"type_detail"`       //流转详细类型 0 未知 1 USDT充值 2银行卡充值 3现货划转合约 4合约划转现货 5提现 6空投支出 7空投收入 8现货支出 9现货收入 10合约支出 11合约收入 12期权支出 13期权收入
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       helpers.TimeNormal `json:"deleted_at"`        //删除时间，为 null 则是没删除
}
