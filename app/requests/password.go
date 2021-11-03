package requests

// Password 用户邮箱
type Password struct {
	Email      string `json:"email" form:"email" validate:"required,email,contains=@"`                   // 邮箱
	Code       string `json:"code" form:"code" validate:"required,numeric"`                              // 验证码
	Password   string `json:"password" form:"password" validate:"required,min=6"`                        // 登录密码
	RePassword string `json:"re_password" form:"re_password" validate:"required,min=6,eqfield=Password"` // 确认登录密码
}
