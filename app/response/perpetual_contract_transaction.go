package response

import (
	"goapi/pkg/helpers"
)

// 永续合约交易订单信息

type PerpetualContractTransaction struct {
	Id              uint               `json:"id"`                //主键id
	UserId          int                `json:"user_id"`           //用户id
	Email           string             `json:"email"`             //用户邮箱
	OrderNumber     string             `json:"order_number"`      //订单号
	CurrencyId      int                `json:"currency_id"`       //币种
	CurrencyName    string             `json:"currency_name"`     //币种名称 例如：BTC/USDT（币种/交易对）
	TradingPairId   int                `json:"trading_pair_id"`   //交易对id
	TradingPairName string             `json:"trading_pair_name"` //交易对名称
	KLineCode       string             `json:"k_line_code"`       //K线图代码
	OrderType       string             `json:"order_type"`        //订单类型：1-限价 2-市价
	LimitPrice      string             `json:"limit_price"`       //当前限价
	TransactionType string             `json:"transaction_type"`  //交易类型：1-开多 2-开空
	EntrustNum      string             `json:"entrust_num"`       //委托量
	EntrustPrice    float64            `json:"entrust_price"`     //委托价格
	ClinchPrice     float64            `json:"clinch_price"`      //成交价格
	EnsureAmount    float64            `json:"ensure_amount"`     //保证金
	HandleFee       float64            `json:"handle_fee"`        //手续费，单位百分比
	Multiple        int                `json:"multiple"`          //倍数值
	Price           float64            `json:"price"`             //交易金额
	Income          float64            `json:"income"`            //最终收益
	Status          int                `json:"status"`            //状态：0 交易中 1 已完成 2 已撤单
	Remark          string             `json:"remark"`            //备注
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       helpers.TimeNormal `json:"deleted_at"`        //删除时间，为 null 则是没删除
}
