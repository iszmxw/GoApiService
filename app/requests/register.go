package requests

// UserEmail 用户邮箱
type UserEmail struct {
	Email string `json:"email" form:"email" validate:"required,email"` // 邮箱
}

// SendEmailRegister 用户邮箱
type SendEmailRegister struct {
	Email      string `json:"email" form:"email" validate:"required,email,contains=@"`          // 邮箱
	ShareCode  string `json:"share_code" form:"share_code"`                                     // 邀请码
	Password   string `json:"password" form:"password" validate:"min=6"`                        // 密码
	RePassword string `json:"re_password" form:"re_password" validate:"min=6,eqfield=Password"` // 确认密码
}

// UserRegister 用户邮箱
type UserRegister struct {
	Email      string `json:"email" form:"email" validate:"required,email,contains=@"`          // 邮箱
	Code       string `json:"code" form:"code" validate:"required,numeric"`                     // 验证码
	ShareCode  string `json:"share_code" form:"share_code"`                                     // 邀请码
	Password   string `json:"password" form:"password" validate:"min=6"`                        // 密码
	RePassword string `json:"re_password" form:"re_password" validate:"min=6,eqfield=Password"` // 确认密码
}
