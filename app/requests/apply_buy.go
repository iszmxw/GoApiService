package requests

// SubmitApplyBuy 提交申购表单
type SubmitApplyBuy struct {
	GetCurrencyId  int `json:"get_currency_id" form:"get_currency_id" validate:"required,numeric"`   // 购买币种id
	GetCurrencyNum int `json:"get_currency_num" form:"get_currency_num" validate:"required,numeric"` // 购买数量
}
