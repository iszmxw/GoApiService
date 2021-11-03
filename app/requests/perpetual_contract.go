package requests

type PerpetualContract struct {
	CurrencyId int `json:"currency_id" form:"currency_id" validate:"required"` //币种id
	Type       int `json:"type" form:"type"`                                   //类型:1.手数,2.倍数
}
