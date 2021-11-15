package requests

// PayPassword 用户邮箱
type PayPassword struct {
	Email         string `json:"email" form:"email" validate:"required,email,contains=@"`                              // 邮箱
	PayPassword   string `json:"pay_password" form:"pay_password" validate:"required,min=6"`                           // 支付密码
	RePayPassword string `json:"re_pay_password" form:"re_pay_password" validate:"required,min=6,eqfield=PayPassword"` // 确认支付密码
}
