package response

import (
	"goapi/pkg/helpers"
)

type PerpetualContract struct {
	Id         int64              `json:"id"`          //主键id
	CurrencyId int                `json:"currency_id"` //币种id
	Type       int8               `json:"type"`        //类型:1.手数,2.倍数
	Value      string             `json:"value"`       //值
	Status     int8               `json:"status"`      //状态:0.禁用,1.启用
	CreatedAt  helpers.TimeNormal `json:"created_at"`  //创建时间
	UpdatedAt  helpers.TimeNormal `json:"updated_at"`  //更新时间
	DeletedAt  helpers.TimeNormal `json:"deleted_at"`  //删除时间，为 null 则是没删除
}
