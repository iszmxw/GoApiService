package models

import (
	"gorm.io/gorm"
	"time"
)

type Currency struct {
	Id                   int           `json:"id"`                     //主键id
	Name                 string         `json:"name"`                   //名称
	TradingPairId        int            `json:"trading_pair_id"`        //交易对
	KLineCode            string         `json:"k_line_code"`            //K线图代码
	DecimalScale         int            `json:"decimal_scale"`          //自有币位数
	Type                 string         `json:"type"`                   //交易显示：（币币交易，永续合约，期权合约）
	Sort                 int8           `json:"sort"`                   //排序
	IsHidden             int8           `json:"is_hidden"`              //是否展示：0-否，1-展示
	FluctuationMin       float64        `json:"fluctuation_min"`        //行情波动值（小）
	FluctuationMax       float64        `json:"fluctuation_max"`        //行情波动值（大）
	FeeCurrencyCurrency  string         `json:"fee_currency_currency"`  //币币交易手续费%
	FeePerpetualContract string         `json:"fee_perpetual_contract"` //永续合约手续费%
	FeeOptionContract    string         `json:"fee_option_contract"`    //期权合约手续费%
	CreatedAt            time.Time      `json:"created_at"`             //创建时间
	UpdatedAt            time.Time      `json:"updated_at"`             //更新时间
	DeletedAt            gorm.DeletedAt `json:"deleted_at"`             //删除时间，为 null 则是没删除
}
