package requests

type PerpetualContract struct {
	CurrencyId int `json:"currency_id" form:"currency_id" validate:"required"` //币种id
}
