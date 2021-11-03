package response

import (
	"goapi/pkg/helpers"
)

type UsersWallet struct {
	Id              uint               `json:"id"`                //主键id
	UserId          int                `json:"user_id"`           //用户id
	TradingPairId   int                `json:"trading_pair_id"`   //交易对id
	TradingPairName string             `json:"trading_pair_name"` //交易对名称
	Address         string             `json:"address"`           //钱包地址
	Status          int                `json:"status"`            //0正常 1锁定
	Available       float64            `json:"available"`         //可用
	Freeze          float64            `json:"freeze"`            //冻结
	TotalAssets     float64            `json:"total_assets"`      //总资产
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       helpers.TimeNormal `json:"deleted_at"`        //删除时间，为 null 则是没删除
}
